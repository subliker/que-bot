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

func NewPlaced(memberCount int) *Queue {
	pn := &personNode{}
	var q = Queue{
		mu:   sync.Mutex{},
		head: pn,
		tail: pn,
		ms:   make(map[telegram.SenderID]struct{}),
	}
	for range memberCount {
		pn := &personNode{}
		q.tail.next = pn
		q.tail = q.tail.next
	}
	return &q
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

// Place returns true if sender was placed and false if place is not available.
//
// If person is already in queue and choose available, it removes from old place and takes new.
func (q *Queue) Place(senderID telegram.SenderID, person telegram.Person, place int) bool {
	cur := q.head
	if place < 0 {
		return false
	}

	// move to place
	for range place + 1 {
		cur = cur.next
		if cur == nil {
			return false
		}
	}

	// if place is not available
	if cur.senderID != 0 {
		return false
	}

	// delete if already in queue
	ncur := q.head
	for ncur.next != nil {
		if ncur.next.senderID == senderID {
			ncur.next.senderID = 0
			ncur.next.person = telegram.Person{}
			break
		}
		ncur = ncur.next
	}

	// set place
	cur.person = person
	cur.senderID = senderID
	q.ms[senderID] = struct{}{}
	return true
}

// PlaceHead returns true if sender was placed and false if there are no places.
func (q *Queue) PlaceHead(senderID telegram.SenderID, person telegram.Person) bool {
	cur := q.head

	// delete if already in queue
	ncur := q.head
	for ncur.next != nil {
		if ncur.next.senderID == senderID {
			ncur.next.senderID = 0
			ncur.next.person = telegram.Person{}
			break
		}
		ncur = ncur.next
	}

	// move to place
	for cur.next != nil {
		if cur.next.senderID == 0 {
			cur.next.person = person
			cur.next.senderID = senderID
			q.ms[senderID] = struct{}{}
			return true
		}
		cur = cur.next
	}
	return false
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

func (q *Queue) ClearPlacedSender(senderID telegram.SenderID) bool {
	// get sender index
	_, ok := q.ms[senderID]
	if !ok {
		return false
	}

	cur := q.head
	for cur.next != nil {
		if cur.next.senderID == senderID {
			cur.next.senderID = 0
			cur.next.person = telegram.Person{}
			delete(q.ms, senderID)
			return true
		}
		cur = cur.next
	}
	return false
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

func (q *Queue) LockedPlaceHeadAndList(senderID telegram.SenderID, person telegram.Person) ([]telegram.Person, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// append
	ok := q.PlaceHead(senderID, person)
	if !ok {
		return nil, false
	}

	// get list
	return q.List(), true
}

func (q *Queue) LockedPlaceAndList(senderID telegram.SenderID, person telegram.Person, place int) ([]telegram.Person, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// append
	ok := q.Place(senderID, person, place)
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

func (q *Queue) LockedClearPlacedSenderAndList(senderID telegram.SenderID) ([]telegram.Person, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// delete
	ok := q.ClearPlacedSender(senderID)
	if !ok {
		return nil, false
	}

	// get list
	return q.List(), true
}
