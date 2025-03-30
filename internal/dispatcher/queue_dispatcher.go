package dispatcher

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/subliker/que-bot/internal/dispatcher/queue"
	"github.com/subliker/que-bot/internal/domain/telegram"
)

type QueueDispatcher interface {
	// Add creates new queue if it doesn't exist
	//
	// Returns ErrQueueAlreadyExists if queue with queue id already exists in
	Add(queueID queue.ID) error
	// SubmitSender submits new sender person with sender id in queue with queue id.
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	//
	// Returns ErrQueueSenderAlreadyExists if sender with sender id already exists in queue
	SubmitSender(queueID queue.ID, senderID telegram.SenderID, person telegram.Person) error
	// List returns ordered slice of telegram submitted persons
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	List(queueID queue.ID) ([]telegram.Person, error)
	// SubmitSenderAndList submits sender person and returns actual person
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	//
	// Returns ErrQueueSenderAlreadyExists if sender with sender id already exists in queue
	SubmitSenderAndList(queueID queue.ID, senderID telegram.SenderID, person telegram.Person) ([]telegram.Person, error)
}

type queueDispatcher struct {
	qs *expirable.LRU[queue.ID, queue.Queue]
}

func NewQueueDispatcher(cfg QueueConfig) QueueDispatcher {
	var qd queueDispatcher

	// making lru queue cache
	qd.qs = expirable.NewLRU(cfg.CacheSize, qd.onCleanup, time.Second*time.Duration(cfg.CacheTTL))

	return &qd
}

func (qd *queueDispatcher) Add(queueID queue.ID) error {
	// try if already exists
	ok := qd.qs.Contains(queueID)
	if ok {
		return ErrQueueAlreadyExists
	}

	// add new queue
	qd.qs.Add(queueID, queue.New())
	return nil
}

func (qd *queueDispatcher) SubmitSender(queueID queue.ID, senderID telegram.SenderID, person telegram.Person) error {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return ErrQueueNotExists
	}

	// submit sender
	ok = q.Append(senderID, person)
	if !ok {
		return ErrQueueSenderAlreadyExists
	}
	return nil
}

func (qd *queueDispatcher) List(queueID queue.ID) ([]telegram.Person, error) {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return nil, ErrQueueNotExists
	}

	return q.List(), nil
}

func (qd *queueDispatcher) SubmitSenderAndList(queueID queue.ID, senderID telegram.SenderID, person telegram.Person) ([]telegram.Person, error) {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return nil, ErrQueueNotExists
	}

	// append and get list with lock
	ls, ok := q.LockedAppendAndList(senderID, person)
	if !ok {
		return nil, ErrQueueSenderAlreadyExists
	}
	return ls, nil
}

func (qb *queueDispatcher) onCleanup(queueID queue.ID, q queue.Queue) {

}
