package queue

import (
	"github.com/google/uuid"
)

type ID string

func GenID(queueName string) ID {
	return ID(uuid.NewString()[:8])
}
