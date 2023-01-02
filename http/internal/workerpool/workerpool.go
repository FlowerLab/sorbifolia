package workerpool

import (
	"errors"
	"log"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"go.x2ox.com/sorbifolia/coarsetime"
	"go.x2ox.com/sorbifolia/http/httperr"
)

// ServeHandler must process tls.Config.NextProto negotiated requests.
type ServeHandler func(c net.Conn) error

// WorkerPool serves incoming connections via a pool of workers
// in FILO order, i.e. the most recently stopped worker will serve the next
// incoming connection.
//
// Such a scheme keeps CPU caches hot (in theory).
type WorkerPool struct {
	// Function for serving server connections.
	// It must leave c unclosed.
	WorkerFunc ServeHandler

	MaxWorkersCount int

	MaxIdleWorkerDuration time.Duration

	lock         sync.Mutex
	workersCount int
	mustStop     bool

	ready []*workerChan

	stopCh chan struct{}

	workerChanPool sync.Pool

	// SetConnState func(net.Conn, ConnState)
	// Logger Logger

	idleConn map[net.Conn]time.Time
}

type workerChan struct {
	lastUseTime time.Time
	ch          chan net.Conn
}

func (wp *WorkerPool) SetConnState(c net.Conn, state ConnState) {
	wp.lock.Lock()
	switch state {
	case StateIdle:
		if wp.idleConn == nil {
			wp.idleConn = make(map[net.Conn]time.Time)
		}
		wp.idleConn[c] = time.Now()
	case StateNew:
		if wp.idleConn == nil {
			wp.idleConn = make(map[net.Conn]time.Time)
		}
		// Count the connection as Idle after 5 seconds.
		// Same as net/http.Server: https://github.com/golang/go/blob/85d7bab91d9a3ed1f76842e4328973ea75efef54/src/net/http/server.go#L2834-L2836
		wp.idleConn[c] = coarsetime.Now().Add(time.Second * 5)

	default:
		delete(wp.idleConn, c)
	}
	wp.lock.Unlock()
}

func (wp *WorkerPool) Start() {
	if wp.stopCh != nil {
		panic("BUG: WorkerPool already started")
	}
	wp.stopCh = make(chan struct{})
	stopCh := wp.stopCh
	wp.workerChanPool.New = func() interface{} {
		return &workerChan{
			ch: make(chan net.Conn, workerChanCap),
		}
	}
	go func() {
		var scratch []*workerChan
		for {
			wp.clean(&scratch)
			select {
			case <-stopCh:
				return
			default:
				time.Sleep(wp.getMaxIdleWorkerDuration())
			}
		}
	}()
}

func (wp *WorkerPool) Stop() {
	if wp.stopCh == nil {
		panic("BUG: WorkerPool wasn't started")
	}
	close(wp.stopCh)
	wp.stopCh = nil

	// Stop all the workers waiting for incoming connections.
	// Do not wait for busy workers - they will stop after
	// serving the connection and noticing wp.mustStop = true.
	wp.lock.Lock()
	ready := wp.ready
	for i := range ready {
		ready[i].ch <- nil
		ready[i] = nil
	}
	wp.ready = ready[:0]
	wp.mustStop = true
	wp.lock.Unlock()
}

func (wp *WorkerPool) getMaxIdleWorkerDuration() time.Duration {
	if wp.MaxIdleWorkerDuration <= 0 {
		return 10 * time.Second
	}
	return wp.MaxIdleWorkerDuration
}

func (wp *WorkerPool) clean(scratch *[]*workerChan) {
	maxIdleWorkerDuration := wp.getMaxIdleWorkerDuration()

	// Clean least recently used workers if they didn't serve connections
	// for more than maxIdleWorkerDuration.
	criticalTime := coarsetime.Now().Add(-maxIdleWorkerDuration)

	wp.lock.Lock()
	ready := wp.ready
	n := len(ready)

	// Use binary-search algorithm to find out the index of the least recently worker which can be cleaned up.
	l, r, mid := 0, n-1, 0
	for l <= r {
		mid = (l + r) / 2
		if criticalTime.After(wp.ready[mid].lastUseTime) {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	i := r
	if i == -1 {
		wp.lock.Unlock()
		return
	}

	*scratch = append((*scratch)[:0], ready[:i+1]...)
	m := copy(ready, ready[i+1:])
	for i = m; i < n; i++ {
		ready[i] = nil
	}
	wp.ready = ready[:m]
	wp.lock.Unlock()

	// Notify obsolete workers to stop.
	// This notification must be outside the wp.lock, since ch.ch
	// may be blocking and may consume a lot of time if many workers
	// are located on non-local CPUs.
	tmp := *scratch
	for i := range tmp {
		tmp[i].ch <- nil
		tmp[i] = nil
	}
}

func (wp *WorkerPool) Serve(c net.Conn) bool {
	ch := wp.getCh()
	if ch == nil {
		return false
	}
	ch.ch <- c
	return true
}

var workerChanCap = func() int {
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}
	return 1
}()

func (wp *WorkerPool) getCh() *workerChan {
	var ch *workerChan
	createWorker := false

	wp.lock.Lock()
	ready := wp.ready
	n := len(ready) - 1
	if n < 0 {
		if wp.workersCount < wp.MaxWorkersCount {
			createWorker = true
			wp.workersCount++
		}
	} else {
		ch = ready[n]
		ready[n] = nil
		wp.ready = ready[:n]
	}
	wp.lock.Unlock()

	if ch == nil {
		if !createWorker {
			return nil
		}
		vch := wp.workerChanPool.Get()
		ch = vch.(*workerChan)
		go func() {
			wp.workerFunc(ch)
			wp.workerChanPool.Put(vch)
		}()
	}
	return ch
}

func (wp *WorkerPool) release(ch *workerChan) bool {
	ch.lastUseTime = coarsetime.Now()
	wp.lock.Lock()
	if wp.mustStop {
		wp.lock.Unlock()
		return false
	}
	wp.ready = append(wp.ready, ch)
	wp.lock.Unlock()
	return true
}

func (wp *WorkerPool) workerFunc(ch *workerChan) {
	var c net.Conn

	var err error
	for c = range ch.ch {
		if c == nil {
			break
		}

		wp.SetConnState(c, StateActive)
		if err = wp.WorkerFunc(c); err != nil && err != httperr.ErrHijacked {
			errStr := err.Error()
			if !(strings.Contains(errStr, "broken pipe") ||
				strings.Contains(errStr, "reset by peer") ||
				strings.Contains(errStr, "request headers: small read buffer") ||
				strings.Contains(errStr, "unexpected EOF") ||
				strings.Contains(errStr, "i/o timeout") ||
				errors.Is(err, httperr.ErrBadTrailer)) {
				log.Printf("error when serving connection %q<->%q: %v", c.LocalAddr(), c.RemoteAddr(), err)
			}
		}
		if err == httperr.ErrHijacked {
			wp.SetConnState(c, StateHijacked)
		} else {
			_ = c.Close()
			wp.SetConnState(c, StateClosed)
		}
		c = nil

		if !wp.release(ch) {
			break
		}
	}

	wp.lock.Lock()
	wp.workersCount--
	wp.lock.Unlock()
}
