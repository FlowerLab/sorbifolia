package httpconfig

import (
	"time"
)

var (
	defaultName   = []byte{'S', 'o', 'r', 'b', 'i', 'f', 'o', 'l', 'i', 'a'}
	defaultConfig = &Config{}
)

const (
	defaultMaxRequestHeaderSize = 4 * 1024
	defaultMaxRequestBodySize   = 4 * 1024 * 1024
	defaultConcurrency          = 256 * 1024

	defaultMaxRequestMethodSize = 7        // Up to 7 if it has not custom methods
	defaultMaxRequestURISize    = 4 * 1024 // Up to 7 if it has not custom methods

	defaultReadTimeout = 60 * time.Second
)

type Config struct {
	// Name is the name of the Server/Client, if not set use defaultName
	// 		Server.Response.Header.Server: Name
	// 		Client.Request.Header.User-Agent: Name
	Name []byte

	MaxRequestMethodSize  int // 最大首行大小
	MaxRequestURISize     int // 最大首行大小
	MaxRequestHeaderSize  int // 最大允许的头大小，包括首行和 \r\n
	MaxRequestBodySize    int // 最大允许的 Body 大小
	StreamRequestBodySize int // 最大允许内存读入的 Body 大小

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	MaxIdleWorkerDuration              time.Duration
	SleepWhenConcurrencyLimitsExceeded time.Duration

	Concurrency int
}

func (c *Config) GetName() []byte { return aObB(dc(c).Name, defaultName) }
func (c *Config) GetMaxRequestMethodSize() int {
	return aObI(dc(c).MaxRequestMethodSize, defaultMaxRequestMethodSize)
}
func (c *Config) GetMaxRequestURISize() int {
	return aObI(dc(c).MaxRequestURISize, defaultMaxRequestURISize)
}
func (c *Config) GetMaxRequestHeaderSize() int {
	return aObI(dc(c).MaxRequestHeaderSize, defaultMaxRequestHeaderSize)
}
func (c *Config) GetMaxRequestBodySize() int {
	return aObI(dc(c).MaxRequestBodySize, defaultMaxRequestBodySize)
}
func (c *Config) GetConcurrency() int           { return aObI(dc(c).Concurrency, defaultConcurrency) }
func (c *Config) GetReadTimeout() time.Duration { return aObD(dc(c).ReadTimeout, defaultReadTimeout) }

func dc(c *Config) *Config {
	if c != nil {
		return c
	}
	return defaultConfig
}

func aObB(a, b []byte) []byte {
	if a != nil {
		return a
	}
	return b
}

func aObI(a, b int) int {
	if a != 0 {
		return a
	}
	return b
}
func aObD(a, b time.Duration) time.Duration {
	if a != 0 {
		return a
	}
	return b
}
