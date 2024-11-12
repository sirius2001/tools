package asyn_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/sirius2001/tools/asyn"
)

func TestAtomicQueque(t *testing.T) {
	q := asyn.NewAtomicQueue()
	wg := sync.WaitGroup{}

	for i := range 100 {
		wg.Add(1)
		go func() {
			q.Enqueue(i)
			wg.Done()
		}()
	}
	wg.Wait()

	for range 100 {
		wg.Add(1)
		go func() {
			res := q.Dequeue()
			fmt.Println(res)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestLockQueque(t *testing.T) {
	q := asyn.NewLockQueue()
	wg := sync.WaitGroup{}

	for i := range 100 {
		wg.Add(1)
		go func() {
			q.Enqueue(i)
			wg.Done()
		}()
	}
	wg.Wait()
	for range 100 {
		wg.Add(1)
		go func() {
			res, exit := q.Dequeue()
			if !exit {
				panic("出现并发问题")
			}
			fmt.Println(res)
			wg.Done()
		}()
	}
	wg.Wait()
}
