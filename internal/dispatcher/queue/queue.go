package queue

import (
	"sync"

	"github.com/subliker/que-bot/internal/domain/telegram"
)

type personNode struct {
	senderID telegram.SenderID
	person   telegram.Person
	next     *personNode
}

type Queue struct {
	ms   map[telegram.SenderID]struct{}
	head *personNode
	tail *personNode
	mu   sync.Mutex
}

func New() *Queue {
	pn := &personNode{}
	return &Queue{
		mu:   sync.Mutex{},
		head: pn,
		tail: pn,
		ms:   make(map[telegram.SenderID]struct{}),
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
	pn := &personNode{
		senderID: senderID,
		person:   person,
	}
	q.tail.next = pn
	q.tail = q.tail.next
	q.ms[senderID] = struct{}{}
	return true
}

func (q *Queue) Delete(senderID telegram.SenderID) bool {
	// get sender index
	_, ok := q.ms[senderID]
	if !ok {
		return false
	}

	// remove sender and move arr elems
	cur := q.head
	for cur.next != nil {
		if cur.next.senderID == senderID {
			ok = true
			if cur.next == q.tail {
				q.tail = cur
			}
			cur.next = cur.next.next
			break
		}
		cur = cur.next
	}
	delete(q.ms, senderID)
	if !ok {
		panic("queue elem wasn't delete")
	}
	return true
}

func (q *Queue) List() []telegram.Person {
	arr := make([]telegram.Person, 0)
	cur := q.head
	for cur.next != nil {
		arr = append(arr, cur.next.person)
		cur = cur.next
	}
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

func (q *Queue) LockedDeleteAndList(senderID telegram.SenderID) ([]telegram.Person, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// delete
	ok := q.Delete(senderID)
	if !ok {
		return nil, false
	}

	// get list
	return q.List(), true
}
