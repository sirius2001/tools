package asyn

import (
	"sync/atomic"
	"unsafe"
)

type QueueNode struct {
	value any
	next  unsafe.Pointer
}

func newNode(v any) *QueueNode {
	return &QueueNode{
		value: v,
		next:  nil,
	}
}

type Queue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func (q *Queue) Enqueue(v any) {
	node := newNode(v)
	for {
		//取出队列的尾部指针
		tail := load(&q.tail)
		//取出尾部指针的next指针
		next := tail.next
		//判定此时的微博的微博指针与队列的微博指针是否相同 如果这里不同 说明其他go操作了队列
		if tail == load(&q.tail) {
			//微博指针的next一般为空
			if next == nil {
				//判断队列的尾部的next指针没有发生了变化 没有则替换新的指针
				if atomic.CompareAndSwapPointer(&tail.next, next, unsafe.Pointer(node)) {
					//判定队列的尾部指针没有发生了变换 没有则将队列的尾部指针替换一下
					atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(node))
					break
				}
			} else {
				//尝试找到队列的微列的微博指针
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), next)
			}
		}
	}
}

func (q *Queue) Dequeue() (res any) {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := head.next

		if head == load(&q.head) {
			//如果队列头部后尾部相同 有空能队列空了
			if head == tail {
				if next == nil {
					return nil // has no value
				}
				//如果突然tail.next不为空了 说明go在尾部添加了一个新元素
				//如果此时队列的微博指针没有发生变化 则将队列尾部指针向后移一位
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), next)
			} else {
				v := load(&next).value
				//如果队列的头没有变换 则将队列的头移出
				if atomic.CompareAndSwapPointer(&q.head, unsafe.Pointer(head), next) {
					return v
				}
			}
		}
	}
}

func NewAtomicQueue() *Queue {
	node := unsafe.Pointer(&QueueNode{})
	return &Queue{
		head: node,
		tail: node,
	}
}

func load(p *unsafe.Pointer) *QueueNode {
	return (*QueueNode)(atomic.LoadPointer(p))
}
