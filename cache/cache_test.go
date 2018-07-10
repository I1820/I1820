package cache

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

type student struct {
	Name   string
	Family string
}

func TestCache1(t *testing.T) {
	rd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	s1 := student{
		Name:   "Parham",
		Family: "Alvani",
	}

	s2 := student{
		Name:   "Navid",
		Family: "Mashayekhi",
	}

	c := &Cache{
		Redis: rd,
		Name:  "students",

		Codec: MsgpackCodec{},
	}

	t.Run("Push 2", func(t *testing.T) {
		assert.NoError(t, c.Push(s1))
		assert.NoError(t, c.Push(s2))
	})

	t.Run("Pop 2", func(t *testing.T) {
		var v student

		assert.NoError(t, c.Pop(&v))
		assert.Equal(t, s2, v)

		assert.NoError(t, c.Pop(&v))
		assert.Equal(t, s1, v)
	})
}
