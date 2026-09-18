package db

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		log.Fatal("REDIS_URL is not set")
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("redis url parse error: %v", err)
	}

	Rdb = redis.NewClient(opt)

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("redis ping error: %v", err)
	}

	log.Println("✅ Connected to Redis")
}

// PollChannel returns the pub/sub channel name used for a given poll's live updates.
func PollChannel(pollID string) string {
	return "poll:" + pollID
}
