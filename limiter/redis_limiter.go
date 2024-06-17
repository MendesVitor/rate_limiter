package limiter

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLimiter struct {
	client         *redis.Client
	ipRateLimit    int
	tokenRateLimit int
	blockDuration  time.Duration
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	ipRateLimit, _ := strconv.Atoi(os.Getenv("IP_RATE_LIMIT"))
	tokenRateLimit, _ := strconv.Atoi(os.Getenv("TOKEN_RATE_LIMIT"))
	blockDuration, _ := strconv.Atoi(os.Getenv("BLOCK_DURATION"))

	return &RedisLimiter{
		client:         client,
		ipRateLimit:    ipRateLimit,
		tokenRateLimit: tokenRateLimit,
		blockDuration:  time.Duration(blockDuration) * time.Second,
	}
}

func (rl *RedisLimiter) Allow(ip string, token string) (bool, error) {
	ctx := context.Background()

	if token != "" {
		tokenKey := "token:" + token
		blockedKey := "blocked:token:" + token

		_, err := rl.client.Get(ctx, blockedKey).Result()
		if err == redis.Nil {
			count, err := rl.client.Incr(ctx, tokenKey).Result()
			if err != nil {
				return false, err
			}
			if count == 1 {
				rl.client.Expire(ctx, tokenKey, time.Second)
			}
			if count > int64(rl.tokenRateLimit) {
				rl.client.Set(ctx, blockedKey, "blocked", rl.blockDuration).Err()
				return false, nil
			}
		} else if err != nil {
			return false, err
		} else {
			return false, nil
		}
	} else {
		ipKey := "ip:" + ip
		blockedKey := "blocked:ip:" + ip

		_, err := rl.client.Get(ctx, blockedKey).Result()
		if err == redis.Nil {
			count, err := rl.client.Incr(ctx, ipKey).Result()
			if err != nil {
				return false, err
			}
			if count == 1 {
				rl.client.Expire(ctx, ipKey, time.Second)
			}
			if count > int64(rl.ipRateLimit) {
				rl.client.Set(ctx, blockedKey, "blocked", rl.blockDuration).Err()
				return false, nil
			}
		} else if err != nil {
			return false, err
		} else {
			return false, nil
		}
	}
	return true, nil
}

func (rl *RedisLimiter) BlockDuration() time.Duration {
	return rl.blockDuration
}
