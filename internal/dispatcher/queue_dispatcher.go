package dispatcher

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/kr/pretty"
	"github.com/subliker/logger"
	"github.com/subliker/que-bot/internal/dispatcher/queue"
	"github.com/subliker/que-bot/internal/domain/telegram"
)

type QueueDispatcher interface {
	// Add creates new queue if it doesn't exist
	//
	// Returns ErrQueueAlreadyExists if queue with queue id already exists in
	Add(queueID queue.ID) error
	// AddPlaced creates new placed queue if it doesn't exist
	//
	// Returns ErrQueueAlreadyExists if queue with queue id already exists in
	AddPlaced(queueID queue.ID, queueMemberCount int) error
	// SubmitSender submits new sender person with sender id in queue with queue id.
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	//
	// Returns ErrQueueSenderAlreadyExists if sender with sender id already exists in queue
	SubmitSender(queueID queue.ID, senderID telegram.SenderID, person telegram.Person) error
	// SubmitPlacedSender submits new sender person with sender id in queue with queue id in memberPlace place.
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	//
	// Returns ErrQueueMemberCountIncorrect if place is not available
	SubmitPlacedSender(queueID queue.ID, senderID telegram.SenderID, person telegram.Person, memberPlace int) error
	// RemoveSender removes sender from queue with queue id.
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	//
	// Returns ErrQueueSenderNotExists if sender with sender id doesn't exist in queue
	RemoveSender(queueID queue.ID, senderID telegram.SenderID) error
	// List returns ordered slice of telegram submitted persons
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	List(queueID queue.ID) ([]telegram.Person, error)
	// SubmitSenderAndList submits sender person and returns actual persons list
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	//
	// Returns ErrQueueSenderAlreadyExists if sender with sender id already exists in queue
	SubmitSenderAndList(queueID queue.ID, senderID telegram.SenderID, person telegram.Person) (persons []telegram.Person, err error)
	// SubmitPlacedSenderAndList submits sender person in memberPlace place and returns actual persons list
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	//
	// Returns ErrQueueMemberCountIncorrect if place is not available
	SubmitPlacedSenderAndList(queueID queue.ID, senderID telegram.SenderID, person telegram.Person, memberPlace int) (persons []telegram.Person, err error)
	// RemoveSenderAndList removes sender from queue and returns actual persons list
	//
	// Returns ErrQueueNotExists if queue with queue id doesn't exist.
	//
	// Returns ErrQueueSenderNotExists if sender with sender id doesn't exist in queue
	RemoveSenderAndList(queueID queue.ID, senderID telegram.SenderID) (persons []telegram.Person, err error)
}

type queueDispatcher struct {
	qs *expirable.LRU[queue.ID, *queue.Queue]

	logger logger.Logger
}

func NewQueueDispatcher(cfg QueueConfig, logger logger.Logger) QueueDispatcher {
	var qd queueDispatcher

	// set logger
	qd.logger = logger.WithFields("layer", "queue_dispatcher")

	// making lru queue cache
	qd.qs = expirable.NewLRU(cfg.CacheSize, qd.onCleanup, time.Second*time.Duration(cfg.CacheTTL))

	qd.logger.Info("queue dispatcher was initialized")
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

	qd.logger.Debugf("queue(%s) was added", queueID)
	return nil
}

func (qd *queueDispatcher) AddPlaced(queueID queue.ID, queueMemberCount int) error {
	// try if already exists
	ok := qd.qs.Contains(queueID)
	if ok {
		return ErrQueueAlreadyExists
	}

	// add new placed queue
	qd.qs.Add(queueID, queue.NewPlaced(queueMemberCount))

	qd.logger.Debugf("placed queue(%s) was added", queueID)
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

	qd.logger.Debugf("sender(%s) was submitted in queue(%s) with data: \n%# v", senderID, queueID, pretty.Formatter(person))
	return nil
}

func (qd *queueDispatcher) SubmitPlacedSender(queueID queue.ID, senderID telegram.SenderID, person telegram.Person, memberPlace int) error {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return ErrQueueNotExists
	}

	// submit sender
	ok = q.Place(senderID, person, memberPlace)
	if !ok {
		return ErrQueueSenderAlreadyExists
	}

	qd.logger.Debugf("placed sender(%s) was submitted in queue(%s) with data: \n%# v", senderID, queueID, pretty.Formatter(person))
	return nil
}

func (qd *queueDispatcher) List(queueID queue.ID) ([]telegram.Person, error) {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return nil, ErrQueueNotExists
	}

	// get list
	lst := q.List()

	qd.logger.Debugf("queue(%s) was listed: \n%# v", queueID, pretty.Formatter(lst))
	return lst, nil
}

func (qd *queueDispatcher) SubmitSenderAndList(queueID queue.ID, senderID telegram.SenderID, person telegram.Person) ([]telegram.Person, error) {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return nil, ErrQueueNotExists
	}

	// append and get list with lock
	lst, ok := q.LockedAppendAndList(senderID, person)
	if !ok {
		return nil, ErrQueueSenderAlreadyExists
	}

	qd.logger.Debugf("queue(%s) was submitted with sender(%s) with data: \n%# v\n and listed: \n%# v", queueID, senderID, pretty.Formatter(person), pretty.Formatter(lst))
	return lst, nil
}

func (qd *queueDispatcher) SubmitPlacedSenderAndList(queueID queue.ID, senderID telegram.SenderID, person telegram.Person, memberPlace int) (persons []telegram.Person, err error) {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return nil, ErrQueueNotExists
	}

	// append and get list with lock
	lst, ok := q.LockedPlaceAndList(senderID, person, memberPlace)
	if !ok {
		return nil, ErrQueueMemberCountIncorrect
	}

	qd.logger.Debugf("placed queue(%s) was submitted with sender(%s) with data: \n%# v\n and listed: \n%# v", queueID, senderID, pretty.Formatter(person), pretty.Formatter(lst))
	return lst, nil
}

func (qd *queueDispatcher) RemoveSender(queueID queue.ID, senderID telegram.SenderID) error {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return ErrQueueNotExists
	}

	// remove sender
	ok = q.Delete(senderID)
	if !ok {
		return ErrQueueSenderNotExists
	}

	qd.logger.Debugf("sender(%s) was removed from queue(%s)", senderID, queueID)
	return nil
}

func (qd *queueDispatcher) RemoveSenderAndList(queueID queue.ID, senderID telegram.SenderID) ([]telegram.Person, error) {
	// get queue
	q, ok := qd.qs.Get(queueID)
	if !ok {
		return nil, ErrQueueNotExists
	}

	// dele and get list with lock
	lst, ok := q.LockedDeleteAndList(senderID)
	if !ok {
		return nil, ErrQueueSenderNotExists
	}

	qd.logger.Debugf("sender(%s) was deleted from queue(%s) and queue listed: \n%# v", senderID, queueID, pretty.Formatter(lst))
	return lst, nil
}

func (qd *queueDispatcher) onCleanup(queueID queue.ID, q *queue.Queue) {

	qd.logger.Debugf("queue(%s) was cleaned up")
}
