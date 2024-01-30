package go_redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var client = redis.NewClient(&redis.Options{
	DB:   0,
	Addr: "localhost:6379",
})

func TestConnection(t *testing.T) {
	assert.NotNil(t, client)

	//err := client.Close()
	//
	//assert.Nil(t, err)
}

var ctx = context.Background()

func TestPing(t *testing.T) {
	result, err := client.Ping(ctx).Result()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", result)
}

func TestString(t *testing.T) {
	client.SetEx(ctx, "name", "test", 3*time.Second)

	result, err := client.Get(ctx, "name").Result()
	assert.Nil(t, err)
	assert.Equal(t, "test", result)

	time.Sleep(5 * time.Second)
	result, err = client.Get(ctx, "name").Result()
	assert.NotNil(t, err)
}
