package dispatcher

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/subliker/que-bot/internal/dispatcher/queue"
	"github.com/subliker/que-bot/internal/domain/telegram"
)

type QueueDispatcher interface {
	// Add creates new queue if it doesn't exist
	Add(queryID telegram.QueryID) error
	// SubmitSender submits new sender with sender id in queue with query id. Returns number of sender in queue.
	//
	// Returns ErrQueueNotExists if queue with query id doesn't exist.
	//
	// Returns ErrQueueSenderAlreadyExists if sender with sender id already exists in queue
	SubmitSender(queryID telegram.QueryID, senderID telegram.SenderID) (num int, err error)
}

type queueDispatcher struct {
	qs *expirable.LRU[telegram.QueryID, queue.Queue]
}

func NewQueueDispatcher(cfg QueueConfig) QueueDispatcher {
	var qd queueDispatcher

	// making lru queue cache
	qd.qs = expirable.NewLRU(cfg.CacheSize, qd.onCleanup, time.Second*time.Duration(cfg.CacheTTL))

	return &qd
}

func (qd *queueDispatcher) Add(queryID telegram.QueryID) error {
	// try if already exists
	ok := qd.qs.Contains(queryID)
	if ok {
		return ErrQueueAlreadyExists
	}

	// add new queue
	qd.qs.Add(queryID, queue.New())
	return nil
}

func (qd *queueDispatcher) SubmitSender(queryID telegram.QueryID, senderID telegram.SenderID) (int, error) {
	// get queue
	q, ok := qd.qs.Get(queryID)
	if !ok {
		return 0, ErrQueueNotExists
	}

	// submit sender
	num, ok := q.Append(senderID)
	if !ok {
		return 0, ErrQueueSenderAlreadyExists
	}
	return num, nil
}

func (qb *queueDispatcher) onCleanup(queryID telegram.QueryID, q queue.Queue) {

}
