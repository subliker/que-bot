package queue

import (
	"sync"

	"github.com/google/uuid"
	"github.com/subliker/que-bot/internal/domain/telegram"
)

type ID uuid.UUID

type Queue struct {
	ms  map[telegram.SenderID]struct{}
	arr *[]telegram.Person
	mu  *sync.Mutex
}

func New() Queue {
	var arr = make([]telegram.Person, 0)
	return Queue{
		mu:  &sync.Mutex{},
		arr: &arr,
		ms:  make(map[telegram.SenderID]struct{}),
	}
}

// Append returns true if sender was added and false if it already exists.
func (q *Queue) Append(senderID telegram.SenderID, person telegram.Person) bool {
	// check if exists
	_, ok := q.ms[senderID]
	if ok {
		return false
	}

	// add sender
	q.ms[senderID] = struct{}{}
	*q.arr = append(*q.arr, person)
	return true
}

func (q *Queue) List() []telegram.Person {
	arr := make([]telegram.Person, len(*q.arr))
	copy(arr, *q.arr)
	return arr
}

func (q *Queue) LockedAppendAndList(senderID telegram.SenderID, person telegram.Person) ([]telegram.Person, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// append
	ok := q.Append(senderID, person)
	if !ok {
		return nil, false
	}

	// get list
	return q.List(), true
}
