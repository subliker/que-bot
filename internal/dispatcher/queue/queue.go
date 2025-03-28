package queue

import "github.com/subliker/que-bot/internal/domain/telegram"

type Queue struct {
	c  int
	ms map[telegram.SenderID]struct{}
}

func New() Queue {
	return Queue{
		c:  0,
		ms: make(map[telegram.SenderID]struct{}),
	}
}

// Append returns true if sender was added and false if it already exists.
//
// It increments member count and returns number of added sender.
func (q *Queue) Append(senderID telegram.SenderID) (int, bool) {
	// check if exists
	_, ok := q.ms[senderID]
	if ok {
		return 0, false
	}

	// add sender
	q.ms[senderID] = struct{}{}
	q.c++
	return q.c, true
}
