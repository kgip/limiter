package limiter

import (
	"testing"
	"time"
)

func TestSlidingWindow_Acquire(t *testing.T) {
	window := NewSlidingWindowLimiter(15, 5, time.Second)
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

func TestLeakyBucket_PushTask(t *testing.T) {
	leakyBucket := NewLeakyBucketLimiter(1, 1)
	for {
		time.Sleep(time.Second)
		for i := 0; i < 5; i++ {
			if leakyBucket.PushTask(func() {
				t.Log("执行任务")
			}) {
				t.Log("请求成功")
			} else {
				t.Log("请求失败")
			}
		}
	}
}

func TestTokenBucketLimiter_Acquire(t *testing.T) {
	limiter := NewTokenBucketLimiter(10, 3)
	for {
		time.Sleep(time.Second)
		for i := 0; i < 10; i++ {
			go func() {
				if limiter.Acquire() {
					t.Log("请求成功")
				} else {
					t.Log("请求失败")
				}
			}()
		}
	}
}
