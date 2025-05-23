package dispatcher

import "errors"

var (
	ErrQueueNotExists            = errors.New("queue doesn't exist")
	ErrQueueSenderNotExists      = errors.New("queue sender doesn't exist")
	ErrQueueAlreadyExists        = errors.New("queue already exists")
	ErrQueueSenderAlreadyExists  = errors.New("queue sender already exists")
	ErrQueueMemberCountIncorrect = errors.New("member place is not available")
	ErrQueuePlacesOver           = errors.New("places are over")
)
