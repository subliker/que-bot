package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subliker/que-bot/internal/domain/telegram"
)

func TestAppend(t *testing.T) {
	t.Run("successfully append", func(t *testing.T) {
		assert := assert.New(t)

		q := New()
		ok := q.Append(telegram.SenderID(123), telegram.Person{
			Username:  "1",
			FirstName: "2",
			LastName:  "3",
		})
		assert.True(ok)

		ok = q.Append(telegram.SenderID(456), telegram.Person{
			Username:  "4",
			FirstName: "5",
			LastName:  "6",
		})
		assert.True(ok)

		assert.Equal([]telegram.Person{
			{
				Username:  "1",
				FirstName: "2",
				LastName:  "3",
			},
			{
				Username:  "4",
				FirstName: "5",
				LastName:  "6",
			},
		}, q.List())
	})
}

func TestDelete(t *testing.T) {
	t.Run("successfully delete all queue", func(t *testing.T) {
		assert := assert.New(t)

		q := New()
		ok := q.Append(telegram.SenderID(123), telegram.Person{
			Username:  "1",
			FirstName: "2",
			LastName:  "3",
		})
		assert.True(ok, "appending doesn't work")
		ok = q.Append(telegram.SenderID(456), telegram.Person{
			Username:  "4",
			FirstName: "5",
			LastName:  "6",
		})
		assert.True(ok, "appending doesn't work")

		ok = q.Delete(telegram.SenderID(456))
		assert.True(ok)

		assert.Equal([]telegram.Person{{
			Username:  "1",
			FirstName: "2",
			LastName:  "3",
		}}, q.List())

		ok = q.Delete(telegram.SenderID(123))
		assert.True(ok)

		assert.Equal([]telegram.Person{}, q.List())
	})
}

func TestPlace(t *testing.T) {
	t.Run("successfully placed person", func(t *testing.T) {
		assert := assert.New(t)

		queue := NewPlaced(4)

		tp := telegram.Person{
			FirstName: "a",
		}
		ok := queue.Place(1, tp, 2)
		assert.True(ok)

		assert.Equal([]telegram.Person{
			{}, {}, tp, {},
		}, queue.List())
	})

	t.Run("not in range", func(t *testing.T) {
		assert := assert.New(t)

		queue := NewPlaced(4)

		tp := telegram.Person{
			FirstName: "a",
		}
		ok := queue.Place(1, tp, 4)
		assert.False(ok)

		assert.Equal([]telegram.Person{
			{}, {}, {}, {},
		}, queue.List())
	})
}
