package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist

	zsetKey := "waitinglist:football"
	res, err := rdb.ZAdd(ctx, zsetKey, redis.Z{
		Score:  float64(time.Now().UnixMicro()),
		Member: 2,
	}).Result()

	if err != nil {
		panic(err)
	}

	fmt.Println("result:", res)

	list, err := rdb.ZRangeWithScores(ctx, zsetKey, 0, time.Now().UnixMicro()).Result()
	if err != nil {
		panic(err)
	}

	for _, item := range list {
		fmt.Printf("member: %v, score: %v\n", item.Member, int64(item.Score))

		mStr, ok := item.Member.(string)
		if ok && mStr == "1" {
			res, err := rdb.ZRem(ctx, zsetKey, item.Member).Result()
			if err != nil {
				panic(err)
			} else {
				fmt.Println("-> res:", res)
			}
		}
	}
}
