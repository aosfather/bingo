package utils

import (
	"container/list"
	"fmt"
	"sync"
)
/**
  队列实现
 */
type Queue struct {
	data *list.List
	lock sync.Mutex
}

func NewQueue() *Queue {
	q := new(Queue)
	q.data = list.New()
	return q
}

func (q *Queue) push(v interface{}) {
	defer 	q.lock.Unlock()
	q.lock.Lock()
	q.data.PushFront(v)
}

func (q *Queue) pop() interface{} {
	defer q.lock.Unlock()
	q.lock.Lock()
	iter := q.data.Back()
	if iter!=nil {


	v := iter.Value
	q.data.Remove(iter)
	return v
	}
	return nil
}

func (q *Queue) dump() {
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		fmt.Println("item:", iter.Value)
	}
}
