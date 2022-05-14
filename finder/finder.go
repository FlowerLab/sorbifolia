package finder

func Async[T any](f func() T) chan T {
	ch := make(chan T)
	go func() {
		ch <- f()
	}()
	return ch
}

func AsyncOR[T any](f func() T) <-chan T {
	return Async(f)
}

func AsyncC[T any](f func() T) chan struct{} {
	ch := make(chan struct{})
	go func() {
		f()
		ch <- struct{}{}
	}()
	return ch
}
