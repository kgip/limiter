package limiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type TokenBucketLimiter struct {
	ready          int32
	tokens         int64
	bucketCapacity int64
	tokensPerSec   int64
	lock           *sync.Mutex
}

func NewTokenBucketLimiter(bucketCapacity, tokensPerSec int64) *TokenBucketLimiter {
	if bucketCapacity <= 0 || tokensPerSec <= 0 {
		panic("invalid parameter")
	}
	return &TokenBucketLimiter{bucketCapacity: bucketCapacity, tokensPerSec: tokensPerSec, lock: &sync.Mutex{}}
}

func (tb *TokenBucketLimiter) Acquire() bool {
	if atomic.CompareAndSwapInt32(&tb.ready, 0, 1) {
		go func() {
			for {
				if atomic.LoadInt64(&tb.tokens) >= tb.bucketCapacity {
					continue
				}
				tb.lock.Lock()
				for i := 0; i < int(tb.tokensPerSec); i++ {
					if tb.tokens < tb.bucketCapacity {
						tb.tokens++
					} else {
						break
					}
				}
				tb.lock.Unlock()
				time.Sleep(time.Second)
			}
		}()
	}
	if atomic.LoadInt64(&tb.tokens) <= 0 {
		return false
	}
	tb.lock.Lock()
	defer tb.lock.Unlock()
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}
