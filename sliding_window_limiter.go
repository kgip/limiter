package limiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type SlidingWindowLimiter struct {
	winTime               time.Duration
	maxRequestCountPerWin int64
	currRequests          int64
	winSize               int64
	lock                  *sync.Mutex
	head                  *subWindow
	tail                  *subWindow
}

type subWindow struct {
	counter int64
	next    *subWindow
}

func NewSlidingWindowLimiter(maxRequestCountPerWin, winSize int64, winTime time.Duration) *SlidingWindowLimiter {
	if maxRequestCountPerWin <= 0 || winSize <= 0 || winTime <= 0 {
		panic("invalid parameter")
	}
	return &SlidingWindowLimiter{maxRequestCountPerWin: maxRequestCountPerWin, winSize: winSize, winTime: winTime, lock: &sync.Mutex{}}
}

func (win *SlidingWindowLimiter) Acquire() bool {
	//第一次请求时创建定时器，开启协程开始
	if win.head == nil {
		win.lock.Lock()
		if win.head == nil {
			//初始化环形队列子窗口
			var curr *subWindow
			for i := 0; i < int(win.winSize); i++ {
				if curr == nil {
					curr = &subWindow{}
					win.head = curr
				} else {
					curr.next = &subWindow{}
					curr = curr.next
				}
			}
			win.tail = curr
			win.tail.next = win.head
			go func() {
				for {
					win.lock.Lock()
					win.currRequests -= win.head.counter
					win.head.counter = 0
					win.tail = win.tail.next
					win.head = win.head.next
					win.lock.Unlock()
					time.Sleep(win.winTime)
				}
			}()
		}
		win.lock.Unlock()
	}
	if atomic.LoadInt64(&win.currRequests) >= win.maxRequestCountPerWin {
		return false
	}
	win.lock.Lock()
	defer win.lock.Unlock()
	if win.currRequests < win.maxRequestCountPerWin {
		win.currRequests++
		win.tail.counter++
		return true
	}
	return false
}
