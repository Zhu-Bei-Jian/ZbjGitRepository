package seque_test

import (
	"testing"
	"time"

	"sanguosha.com/sgs_herox/gameutil/seque"
)

func TestQueue(t *testing.T) {
	ch := make(chan func(), 1)

	go func() {
		q := new(seque.Queue)
		var counter int
		ch <- func() {
			q.Push(func(wg seque.WaitGroup) {
				wg.Add(1)
				time.AfterFunc(30*time.Microsecond, func() {
					ch <- func() {
						defer wg.Done()
						counter++
						if counter != 1 {
							t.Error(counter)
						}
					}
				})
			})
			q.Push(func(wg seque.WaitGroup) {
				wg.Add(1)
				time.AfterFunc(20*time.Microsecond, func() {
					ch <- func() {
						defer wg.Done()
						counter++
						if counter != 2 {
							t.Error(counter)
						}
					}
				})
			})
			q.Push(func(wg seque.WaitGroup) {
				wg.Add(1)
				time.AfterFunc(10*time.Microsecond, func() {
					ch <- func() {
						defer wg.Done()
						counter++
						if counter != 3 {
							t.Error(counter)
						}
					}
				})
			})
			q.Push(func(_ seque.WaitGroup) {
				if counter != 3 {
					t.Error(counter)
				}
				close(ch)
			})
		}
	}()

	for f := range ch {
		f()
	}
}
