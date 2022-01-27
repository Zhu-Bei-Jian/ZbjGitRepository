package worker

import (
	"sync/atomic"
	"time"
)

// Worker 代表一个工作协程
// ioservice.IOService also is a Worker
// 单独成立一个接口的原因是: Post接口的使用者不需要也不应该调用 ioservice.IOService 的其他接口, 这样可使程序的语义更明确
type Worker interface {
	Post(func())
}

// AfterPost 在指定时间之后, 在 worker 中执行 f.
// 返回取消函数.
func AfterPost(worker Worker, d time.Duration, f func()) (cancel func()) {
	// 由于异步回调的特性, 必须要保证, 除非在调用 cancel 前, 回调函数已经被执行, 否则在调用 cancel 函数后, 回调函数必须被丢弃(不能够再执行).
	var stop int32
	t := time.AfterFunc(d, func() {
		worker.Post(func() {
			if atomic.CompareAndSwapInt32(&stop, 0, 1) {
				f()
			}
		})
	})
	return func() {
		if atomic.CompareAndSwapInt32(&stop, 0, 1) {
			t.Stop()
		}
	}
}
