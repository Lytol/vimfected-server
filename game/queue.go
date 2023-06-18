package game

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

func (q *Queue) Peek() (Command, bool) {
	if q.ll.Len() == 0 {
		return Command{}, false
	}
	return q.ll.Front().Value.(Command), true
}

func (q *Queue) Unshift() (Command, bool) {
	if q.ll.Len() == 0 {
		return Command{}, false
	}
	return q.ll.Remove(q.ll.Front()).(Command), true
}

func (q *Queue) Each(fn func(Command) bool) {
	for e := q.ll.Front(); e != nil; e = e.Next() {
		if fn(e.Value.(Command)) {
			q.ll.Remove(e)
		}
	}
}

// Clears all previous commands by the same player
func (q *Queue) ClearUntil(cmd Command) {
	e := q.ll.Front()

	for {
		next := e.Next()

		if e == nil || e.Value.(Command).Timestamp == cmd.Timestamp {
			return
		}

		q.ll.Remove(e)
		e = next
	}
}
