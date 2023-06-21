package test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestRedis(t *testing.T) {
	//r := models.RDB
	//r.Set(ctx, "cao", "jin", 5*time.Minute)
}
