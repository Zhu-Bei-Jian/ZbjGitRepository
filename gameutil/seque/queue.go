package seque

import "sanguosha.com/sgs_herox/gameutil"

// WaitGroup 显示控制异步回调任务何时结束
type WaitGroup interface {
	Add(delta int)
	Done()
}

type waitGroup struct {
	cnt int
	q   *Queue

	isPop bool
}

func (w *waitGroup) Add(delta int) {
	w.cnt += delta
}
func (w *waitGroup) Done() {
	w.cnt--
	if w.cnt == 0 {
		w.isPop = true
		w.q.pop()
		w.q.do()
	} else if w.cnt < 0 {
		panic("illegal action")
	}
}

// Queue 执行队列, 用于异步回调过程中消息的时许问题处理
type Queue struct {
	que []func(WaitGroup)
}

// Push 将将要执行的函数 push 到执行队列中,
// 如果当前没有正在执行的任务, 则立即执行当前的 f, 否则等待之前的任务完成后再执行 f
// WaitGroup 使你能够手动控制合适才算完成一次任务执行(f 中可以再次发起异步任务)
func (q *Queue) Push(f func(WaitGroup)) {
	q.que = append(q.que, f)
	if len(q.que) == 1 {
		q.do()
	}
}

func (q *Queue) do() {
	for len(q.que) > 0 {
		wg := new(waitGroup)
		wg.q = q
		gameutil.SafeCall(func() {
			q.que[0](wg)
		})
		if wg.cnt > 0 {
			break
		}
		//可能在执行函数的时候使用done将自己pop过
		if wg.isPop {
			break
		}
		q.pop()
	}
}

func (q *Queue) pop() {
	copy(q.que, q.que[1:])
	q.que = q.que[:len(q.que)-1]
}
