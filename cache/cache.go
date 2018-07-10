package cache

import (
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

// Codec represents encode/decode methods
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

// MsgpackCodec is MessagePack implementation of codec
type MsgpackCodec struct{}

// Marshal returns the MessagePack encoding of v.
func (MsgpackCodec) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)

}

// Unmarshal decodes the MessagePack-encoded data and stores the result in the value pointed to by v.
func (MsgpackCodec) Unmarshal(b []byte, v interface{}) error {
	return msgpack.Unmarshal(b, v)

}

// Cache represents codec (for automatically encode/decode data)
// and queue cache based on redis.
type Cache struct {
	Codec Codec
	Name  string

	Redis rediser
}

type rediser interface {
	LPush(key string, values ...interface{}) *redis.IntCmd
	LPop(key string) *redis.StringCmd
}

// Push new value into end of queue
func (c *Cache) Push(v interface{}) error {
	b, err := c.Codec.Marshal(v)
	if err != nil {
		return err
	}

	if _, err := c.Redis.LPush(c.Name, b).Result(); err != nil {
		return err
	}

	return nil
}

// Pop value from head of queue
func (c *Cache) Pop(v interface{}) error {
	b, err := c.Redis.LPop(c.Name).Bytes()
	if err != nil {
		return err
	}

	if err := c.Codec.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}
