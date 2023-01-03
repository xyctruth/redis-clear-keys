package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var (
	l             sync.Mutex
	totalCount    = 0
	totalDelCount = 0
	redisHost     = ""
	redisPassword = ""
	rdb           *redis.Client
)

func main() {
	flag.StringVar(&redisHost, "redis-host", "", "")
	flag.StringVar(&redisPassword, "redis-password", "", "")
	flag.Parse()
	if redisHost == "" || redisPassword == "" {
		fmt.Println("redis-host or redis-password can not be empty")
		return
	}

	err := initRedis()
	if err != nil {
		fmt.Println("init redis failed err :", err)
		return
	}

	ch := make(chan []string, 1000)
	go scanKeys(ch)
	for i := 0; i <= 100; i++ {
		go consumeKeys(ch)
	}
	select {}
}

func scanKeys(ch chan []string) {
	var err error
	var cursor, tempCursor uint64
	for {
		var keys []string
		keys, tempCursor, err = rdb.Scan(context.Background(), cursor, "*", 1000).Result()
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			cursor = tempCursor
		}
		ch <- keys
		if cursor == 0 {
			close(ch)
			fmt.Println("scan done")
			break
		}
	}
}

func consumeKeys(ch chan []string) {
	ctx := context.Background()
	for keys := range ch {
		count := 0
		delCount := 0
		ttlPipe := rdb.Pipeline()
		delPipe := rdb.Pipeline()
		for _, key := range keys {
			ttlPipe.TTL(ctx, key)
		}
		ttlCmds, err := ttlPipe.Exec(ctx)
		for _, cmd := range ttlCmds {
			if result, ok := cmd.(*redis.DurationCmd); ok {
				ttl, err := result.Result()
				if err != nil {
					fmt.Println(err)
				}
				if ttl == -1 {
					key := cmd.Args()[1].(string)
					delPipe.Del(context.Background(), key)
					delCount++
				}
			}
		}
		_, err = delPipe.Exec(ctx)
		if err != nil {
			fmt.Println(err)
		}
		count = count + len(keys)
		counter(count, delCount)
	}
}

func counter(count, delCount int) {
	l.Lock()
	defer l.Unlock()
	totalDelCount = totalDelCount + delCount
	totalCount = totalCount + count
	fmt.Printf("total scan num:%d ,total clear num: %d \n", totalCount, totalDelCount)
}

func initRedis() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		PoolSize: 200,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("ping redis failed err:", err)
		return err
	}
	return nil
}
