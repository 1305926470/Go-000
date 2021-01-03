package limiter

import (
	"sync"
	"time"
)

type bucket struct {
	val int64
}

type SlideWindow struct {
	windowSize int64

	bucketMutex *sync.RWMutex
	Buckets     map[int64]*bucket
	bucketCache []*bucket
}

// 窗口计数汇总
func (sw *SlideWindow) Reduce(now time.Time) int64 {
	sw.bucketMutex.RLock()
	defer sw.bucketMutex.RUnlock()
	var sum int64 = 0
	for t, bucket := range sw.Buckets {
		if t > now.Unix()-sw.windowSize {
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
	sw.IncN(1)
}

// 计数
func (sw *SlideWindow) IncN(i int64) {
	sw.bucketMutex.Lock()
	bucket := sw.getCurrentBucket()
	bucket.val += i
	sw.bucketMutex.Unlock()
	sw.removeBucket()
}

// 返回滑动窗口最大bucket的统计数量
func (sw *SlideWindow) Max(now time.Time) int64 {
	var max int64
	sw.bucketMutex.RLock()
	defer sw.bucketMutex.RUnlock()

	for timestamp, bucket := range sw.Buckets {
		if timestamp >= now.Unix()-sw.windowSize {
			if bucket.val > max {
				max = bucket.val
			}
		}
	}
	return max
}

// 计算 bucket 平均计数
func (sw *SlideWindow) Avg(now time.Time) int64 {
	return sw.Reduce(now) / sw.windowSize
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

// 系统统计数据详情
func (sw *SlideWindow) Stats() map[int64]*bucket {
	return sw.Buckets
}

// 按秒级维度划分， 1秒 一个 bucket
// size 为 bucket 数量
func NewWindow(size int64) SlideWindow {
	if size <= 0 {
		panic("The size must be greater than 0")
	}
	return SlideWindow{
		windowSize:  size,
		bucketMutex: &sync.RWMutex{},
		Buckets:     make(map[int64]*bucket, size),
	}
}




//--------------------------------------------------------------------------------
// 统计请求中成功，失败，拒绝，超时等情况
// 在实际应用中，可统计更多等类型
type Collector struct {
	mu        *sync.RWMutex
	successes *SlideWindow
	failures  *SlideWindow
	rejects   *SlideWindow
	timeout   *SlideWindow
}

type Metric struct {
	Successes int64
	Failures  int64
	Rejects   int64
	Timeouts  int64
}

func (c Collector) Update(mtr Metric) {
	c.successes.IncN(mtr.Successes)
	c.failures.IncN(mtr.Failures)
	c.rejects.IncN(mtr.Rejects)
	c.timeout.IncN(mtr.Timeouts)
}

