package limiter

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucketLimiter struct {
	taskQueue       chan func()
	queueSize       int
	execTasksPerSec int
	lock            *sync.Mutex
}

func NewLeakyBucketLimiter(queueSize, execTasksPerSec int) *LeakyBucketLimiter {
	if queueSize <= 0 || execTasksPerSec <= 0 {
		panic("invalid parameter")
	}
	return &LeakyBucketLimiter{queueSize: queueSize, execTasksPerSec: execTasksPerSec, lock: &sync.Mutex{}}
}

func (lb *LeakyBucketLimiter) PushTask(newTask func()) bool {
	if lb.taskQueue == nil {
		lb.lock.Lock()
		if lb.taskQueue == nil {
			lb.taskQueue = make(chan func(), lb.queueSize)
			go func() {
				for {
					for i := 0; i < lb.execTasksPerSec; i++ {
						select {
						case task := <-lb.taskQueue:
							go func() {
								defer func() {
									if err := recover(); err != nil {
										fmt.Println(err)
									}
								}()
								task()
							}()
						default:
							continue
						}
					}
					time.Sleep(time.Second)
				}
			}()
		}
		lb.lock.Unlock()
	}
	select {
	case lb.taskQueue <- newTask:
		return true
	default:
		return false
	}
}
