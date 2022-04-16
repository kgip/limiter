package limiter

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	winTime               time.Duration
	maxRequestCountPerWin int
	currRequests          int
	winSize               int
	lock                  *sync.Mutex
	head                  *subWindow
	tail                  *subWindow
}

type subWindow struct {
	counter int
	next    *subWindow
}

func NewSlidingWindow(maxRequestCountPerWin, winSize int, winTime time.Duration) *SlidingWindow {
	if maxRequestCountPerWin <= 0 || winSize <= 0 || winTime <= 0 {
		panic("invalid parameter")
	}
	return &SlidingWindow{maxRequestCountPerWin: maxRequestCountPerWin, winSize: winSize, winTime: winTime, lock: &sync.Mutex{}}
}

func (win *SlidingWindow) Acquire() bool {
	//第一次请求时创建定时器，开启协程开始
	if win.head == nil {
		win.lock.Lock()
		if win.head == nil {
			//初始化环形队列子窗口
			var curr *subWindow
			for i := 0; i < win.winSize; i++ {
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
					time.Sleep(win.winTime)
					win.lock.Lock()
					win.currRequests -= win.head.counter
					win.head.counter = 0
					win.tail = win.tail.next
					win.head = win.head.next
					win.lock.Unlock()
				}
			}()
		}
		win.lock.Unlock()
	}
	win.lock.Lock()
	defer win.lock.Unlock()
	if win.currRequests+1 <= win.maxRequestCountPerWin {
		win.currRequests++
		win.tail.counter++
		return true
	}
	return false
}
