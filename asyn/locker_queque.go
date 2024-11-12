package asyn

import (
	"sync"
)

// 定义队列节点结构
type Node struct {
	value any
	next  *Node
}

// 加锁队列结构
type LockQueue struct {
	head  *Node      // 队列头节点
	tail  *Node      // 队列尾节点
	mutex sync.Mutex // 用于锁定队列的互斥锁
}

// 初始化一个新的加锁队列
func NewLockQueue() *LockQueue {
	node := &Node{} // 哨兵节点，方便操作
	return &LockQueue{
		head: node,
		tail: node,
	}
}

// 入队操作：将元素加入队列尾部
func (q *LockQueue) Enqueue(value any) {
	q.mutex.Lock()         // 加锁
	defer q.mutex.Unlock() // 确保函数结束后解锁

	// 创建新节点并将其添加到队列尾部
	newNode := &Node{value: value}
	q.tail.next = newNode
	q.tail = newNode
}

// 出队操作：从队列头部取出元素
func (q *LockQueue) Dequeue() (any, bool) {
	q.mutex.Lock()         // 加锁
	defer q.mutex.Unlock() // 确保函数结束后解锁

	// 检查队列是否为空
	if q.head.next == nil {
		return nil, false
	}

	// 取出头部节点的下一个节点
	removedNode := q.head.next
	q.head.next = removedNode.next

	// 如果移除的是最后一个节点，更新 tail 指针
	if q.head.next == nil {
		q.tail = q.head
	}

	return removedNode.value, true
}
