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

func TestList(t *testing.T) {
	client.RPush(ctx, "name", "Eko")
	client.RPush(ctx, "name", "Kurniawan")
	client.RPush(ctx, "name", "Khannedy")

	assert.Equal(t, "Eko", client.LPop(ctx, "name").Val())
	assert.Equal(t, "Kurniawan", client.LPop(ctx, "name").Val())
	assert.Equal(t, "Khannedy", client.LPop(ctx, "name").Val())
}

func TestSet(t *testing.T) {
	client.SAdd(ctx, "category", "fashion")
	client.SAdd(ctx, "category", "fashion")
	client.SAdd(ctx, "category", "electronic")
	client.SAdd(ctx, "category", "electronic")

	assert.Equal(t, int64(2), client.SCard(ctx, "category").Val())
	assert.Equal(t, []string{"fashion", "electronic"}, client.SMembers(ctx, "category").Val())
}

func TestSortedSet(t *testing.T) {
	client.ZAdd(ctx, "scores", redis.Z{Score: 100, Member: "Eko"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 85, Member: "Jhon"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 95, Member: "Santy"})

	assert.Equal(t, []string{"Jhon", "Santy", "Eko"}, client.ZRange(ctx, "scores", 0, -1).Val())
	assert.Equal(t, "Eko", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Santy", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Jhon", client.ZPopMax(ctx, "scores").Val()[0].Member)
}

func TestHash(t *testing.T) {
	client.HSet(ctx, "user:1", "id", "1")
	client.HSet(ctx, "user:1", "name", "john")
	client.HSet(ctx, "user:1", "email", "example@gmail.com")

	user := client.HGetAll(ctx, "user:1").Val()

	assert.Equal(t, "1", user["id"])
	assert.Equal(t, "john", user["name"])
	assert.Equal(t, "example@gmail.com", user["email"])
}

func TestGeoLocation(t *testing.T) {
	client.GeoAdd(ctx, "seller", &redis.GeoLocation{
		Name:      "Toko A",
		Latitude:  -6.275112,
		Longitude: 106.813133,
	})

	client.GeoAdd(ctx, "seller", &redis.GeoLocation{
		Name:      "Toko B",
		Latitude:  -6.278195,
		Longitude: 106.809239,
	})

	distance := client.GeoDist(ctx, "seller", "Toko A", "Toko B", "km").Val()
	assert.Equal(t, 0.5504, distance)

	sellers := client.GeoSearch(ctx, "seller", &redis.GeoSearchQuery{
		Longitude:  106.810305,
		Latitude:   -6.275644,
		Radius:     5,
		RadiusUnit: "km",
	}).Val()

	assert.Equal(t, []string{"Toko B", "Toko A"}, sellers)

}

func TestHyperLogLog(t *testing.T) {
	client.PFAdd(ctx, "visitors", "john", "dimas", "ruly", "budi")
	client.PFAdd(ctx, "visitors", "john", "joko", "ruly", "joni")
	client.PFAdd(ctx, "visitors", "jaki", "akagami", "luffy", "joni")

	total := client.PFCount(ctx, "visitors").Val()

	assert.Equal(t, int64(9), total)
}

func TestPipeline(t *testing.T) {
	_, err := client.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.SetEx(ctx, "name", "john doe", 5*time.Second)
		pipeliner.SetEx(ctx, "city", "Indonesia", 5*time.Second)

		return nil
	})

	assert.Nil(t, err)

	assert.Equal(t, "john doe", client.Get(ctx, "name").Val())
	assert.Equal(t, "Indonesia", client.Get(ctx, "city").Val())
}
