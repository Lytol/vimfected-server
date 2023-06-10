package commands

import (
	"container/list"
)

type Queue struct {
	ll *list.List
}

func NewQueue() *Queue {
	return &Queue{
		ll: list.New(),
	}
}

func (q *Queue) Shift(cmd Command) {
	q.ll.PushBack(cmd)
}

func (q *Queue) Unshift() (Command, bool) {
	if q.ll.Len() == 0 {
		return Command{}, false
	}
	return q.ll.Remove(q.ll.Front()).(Command), true
}
