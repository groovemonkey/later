package main

import (
	// "context"
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

func main() {
	fmt.Println("Testing Golang Redis")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)

}
