package limiter

import (
	"testing"
	"time"
)

func TestSlidingWindow_Acquire(t *testing.T) {
	window := NewSlidingWindow(15, 5, time.Second)
	for {
		time.Sleep(time.Second)
		for i := 0; i < 5; i++ {
			go func() {
				if window.Acquire() {
					t.Log("请求成功")
				} else {
					t.Log("请求失败")
				}
			}()
		}
	}
}
