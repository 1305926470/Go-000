package limiter

import (
	"sync"
	"time"
)

type bucket struct {
	val int64
}

type SlideWindow struct {
	windowSize  int64

	bucketMutex *sync.RWMutex
	Buckets     map[int64]*bucket
	bucketCache []*bucket
}

// 窗口计数汇总
func (sw *SlideWindow) Reduce() int64 {
	sw.bucketMutex.RLock()
	defer sw.bucketMutex.RUnlock()
	var sum int64 = 0
	now := time.Now().Unix()
	for t, bucket := range sw.Buckets {
		if t > now - sw.windowSize {
			sum += bucket.val
		}
	}
	return sum
}

// 移除旧的 bucket
func (sw *SlideWindow) removeBucket() {
	critical := time.Now().Unix() - sw.windowSize
	for t, bucket := range sw.Buckets {
		if t <= critical {
			sw.bucketCache = append(sw.bucketCache, bucket)
			delete(sw.Buckets, t)
		}
	}
}

// 计数
func (sw *SlideWindow) Inc() {
	sw.bucketMutex.Lock()
	bucket := sw.getCurrentBucket()
	bucket.val++
	sw.bucketMutex.Unlock()
	sw.removeBucket()
}

// 获取当前 bucket
func (sw *SlideWindow) getCurrentBucket() *bucket {
	now := time.Now().Unix()
	if b, ok := sw.Buckets[now]; ok {
		return b
	}

	if l := len(sw.bucketCache); l > 0 {
		b := sw.bucketCache[l-1]
		sw.bucketCache = sw.bucketCache[:l-1]
		return b
	}

	b := new(bucket)
	sw.Buckets[now] = b
	return b
}

func (sw *SlideWindow) Stats() map[int64]*bucket {
	return sw.Buckets
}

// 按秒级纬度划分， 1秒 一个 bucket
// size 为 bucket 数量
func NewWindow(size int64) SlideWindow {
	return SlideWindow{
		windowSize: size,
		bucketMutex:  &sync.RWMutex{},
		Buckets: make(map[int64]*bucket, size),
	}
}





