package queue

//参考： https://www.cnblogs.com/yourstars/p/15186558.html
// 继续参考 https://github.com/yireyun/go-queue

import (
	"sync/atomic"
	"unsafe"
)

// 定义无锁队列结构
type LKQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node struct {
	value interface{}
	next  unsafe.Pointer
}

// 新建队列，返回一个空队列
func NewLKQueue() *LKQueue {
	n := unsafe.Pointer(&node{})
	return &LKQueue{head: n, tail: n}
}

// 插入，将给定的值v放在队列的尾部
func (q *LKQueue) Enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) {
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n) // 排队完成, 尝试将tail移到插入的节点
					return
				}
			} else { // tail没有指向最后一个节点
				// 将Tail移到下一个节点
				cas(&q.tail, tail, next)
			}
		}
	}
}

// 移除，删除并返回队列头部的值,如果队列为空，则返回nil
func (q *LKQueue) Dequeue() interface{} {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) {
			if head == tail {
				if next == nil {
					return nil
				}

				cas(&q.tail, tail, next)
			} else {
				// 在CAS之前读取值，否则另一个出队列可能释放下一个节点
				v := next.value
				if cas(&q.head, head, next) {
					return v
				}
			}
		}
	}
}

func load(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

// CAS算法
func cas(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}
