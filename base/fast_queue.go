package base

import (
	"sync"
)

type ILock interface {
	LockEx()
	UnLockEx()
}

type Lock struct {
	mutex sync.Mutex
}

func (l *Lock) LockEx() {
	l.mutex.Lock()
}

func (l *Lock) UnLockEx() {
	l.mutex.Unlock()
}

type NoLock struct {
}

func (l *NoLock) LockEx() {

}

func (l *NoLock) UnLockEx() {

}

type IFastQueue interface {
}

type node[T any] struct {
	element T
	next    *node[T]
}

type FastQueue[T any] struct {
	last  *node[T]
	first *node[T]
	size  int
	lock  ILock
	_nil  T
}

func CreateFastQueue[T any](is bool) *FastQueue[T] {
	var fq FastQueue[T]
	fq.size = 0
	fq.first = nil
	fq.last = nil

	if is {
		fq.lock = new(Lock)
	} else {
		fq.lock = new(NoLock)
	}

	return &fq
}

func (f *FastQueue[T]) Clear() {
	for f.last != nil {
		f.Pop()
	}

	f.size = 0
}

func (f *FastQueue[T]) Push(v T) {
	f.Lock()
	n := new(node[T])
	if f.last != nil {
		f.last.next = n
	} else {
		f.first = n
	}

	f.last = n
	n.next = nil
	n.element = v
	f.size++

	f.UnLock()
}

// Pop 队列先进先出
func (f *FastQueue[T]) Pop() T {
	f.Lock()
	if f.first == nil {
		f.UnLock()
		return f._nil
	}

	ret := f.first.element
	f.first = f.first.next
	if f.first == nil {
		f.last = nil
	}

	f.size--

	f.UnLock()

	return ret
}

func (f *FastQueue[T]) Front() T {
	f.Lock()
	if f.first == nil {
		f.UnLock()
		return f._nil
	}

	ret := f.first.element

	f.UnLock()

	return ret
}

func (f *FastQueue[T]) PopFront() {
	f.Lock()
	if f.first == nil {
		f.UnLock()
		return
	}

	f.first = f.first.next
	if f.first == nil {
		f.last = nil
	}

	f.size--

	f.UnLock()
}

func (f *FastQueue[T]) IsEmpty() bool {
	f.Lock()
	ret := (f.size != 0)
	f.UnLock()
	return ret
}

func (f *FastQueue[T]) Lock() {
	f.lock.LockEx()
}

func (f *FastQueue[T]) UnLock() {
	f.lock.UnLockEx()
}

func (f *FastQueue[T]) GetFirst() *node[T] {
	return f.first
}

func (f *FastQueue[T]) GetLast() *node[T] {
	return f.last
}

func (f *FastQueue[T]) SetFirst(n *node[T]) {
	f.first = n
}

func (f *FastQueue[T]) SetLast(n *node[T]) {
	f.last = n
}

func (f *FastQueue[T]) Copy(t *FastQueue[T]) {
	t.Lock()
	ft := t.first
	lt := t.last
	se := t.size

	t.SetFirst(nil)
	t.SetLast(nil)
	t.size = 0
	t.UnLock()

	f.Lock()
	f.first = ft
	f.last = lt
	f.size = se
	f.UnLock()

}
