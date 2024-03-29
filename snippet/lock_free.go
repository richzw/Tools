package snippet

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

// https://github.com/golang-design/lockfree/blob/master/queue.go

// LKQueue - lock free queue
// LKQueue is a lock-free unbounded queue.
type LKQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}
type node struct {
	value interface{}
	next  unsafe.Pointer
}

// NewLKQueue returns an empty queue.
func NewLKQueue() *LKQueue {
	n := unsafe.Pointer(&node{})
	return &LKQueue{head: n, tail: n}
}

// Enqueue puts the given value v at the tail of the queue.
func (q *LKQueue) Enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) { // are tail and next consistent?
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n) // Enqueue is done.  try to swing tail to the inserted node
					return
				}
			} else { // tail was not pointing to the last node
				// try to swing Tail to the next node
				cas(&q.tail, tail, next)
			}
		}
		runtime.Gosched() // https://stackoverflow.com/questions/12291949/why-my-golang-lock-free-queue-always-stuck-there
	}
}

// Dequeue removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *LKQueue) Dequeue() interface{} {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) { // are head, tail, and next consistent?
			if head == tail { // is queue empty or tail falling behind?
				if next == nil { // is queue empty?
					return nil
				}
				// tail is falling behind.  try to advance it
				cas(&q.tail, tail, next)
			} else {
				// read value before CAS otherwise another dequeue might free the next node
				v := next.value
				if cas(&q.head, head, next) {
					return v // Dequeue is done.  return
				}
			}
		}
		runtime.Gosched()
	}
}
func load(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}
func cas(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}

// CQueue - two lock queue
// CQueue is a concurrent unbounded queue which uses two-Lock concurrent queue qlgorithm.
type CQueue struct {
	head  *cnode
	tail  *cnode
	hlock sync.Mutex
	tlock sync.Mutex
}
type cnode struct {
	value interface{}
	next  *cnode
}

// NewCQueue returns an empty CQueue.
func NewCQueue() *CQueue {
	n := &cnode{}
	return &CQueue{head: n, tail: n}
}

// Enqueue puts the given value v at the tail of the queue.
func (q *CQueue) Enqueue(v interface{}) {
	n := &cnode{value: v}
	q.tlock.Lock()
	q.tail.next = n // Link node at the end of the linked list
	q.tail = n      // Swing Tail to node
	q.tlock.Unlock()
}

// Dequeue removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *CQueue) Dequeue() interface{} {
	q.hlock.Lock()
	n := q.head
	newHead := n.next
	if newHead == nil {
		q.hlock.Unlock()
		return nil
	}
	v := newHead.value
	newHead.value = nil
	q.head = newHead
	q.hlock.Unlock()
	return v
}
