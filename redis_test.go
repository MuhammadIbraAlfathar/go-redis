package go_redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

var client = redis.NewClient(&redis.Options{
	DB:   0,
	Addr: "localhost:6379",
})

func TestConnection(t *testing.T) {
	assert.NotNil(t, client)

	err := client.Close()

	assert.Nil(t, err)
}
