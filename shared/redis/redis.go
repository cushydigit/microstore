package redis

import (
	"context"
	"log"

	rds "github.com/redis/go-redis/v9"
)

var (
	Client *rds.Client
	Ctx    = context.Background()
)

func Init(Addr string) {
	Client = rds.NewClient(&rds.Options{
		Addr:     Addr,
		Password: "",
		DB:       0,
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
}
