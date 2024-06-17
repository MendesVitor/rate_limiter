package limiter

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupMockRedis() (*miniredis.Miniredis, *redis.Client, error) {
	mockRedis, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
	})
	return mockRedis, client, nil
}

func TestRedisLimiter_AllowIP(t *testing.T) {
	mockRedis, client, err := setupMockRedis()
	assert.NoError(t, err)
	defer mockRedis.Close()
	defer client.Close()

	ipRateLimit := 5
	blockDuration := 2
	rl := NewRedisLimiter(client)
	rl.ipRateLimit = ipRateLimit
	rl.blockDuration = time.Duration(blockDuration) * time.Second

	for i := 0; i < ipRateLimit; i++ {
		allowed, err := rl.Allow("192.168.1.1", "")
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rl.Allow("192.168.1.1", "")
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestRedisLimiter_AllowToken(t *testing.T) {
	mockRedis, client, err := setupMockRedis()
	assert.NoError(t, err)
	defer mockRedis.Close()
	defer client.Close()

	tokenRateLimit := 10
	blockDuration := 2
	rl := NewRedisLimiter(client)
	rl.tokenRateLimit = tokenRateLimit
	rl.blockDuration = time.Duration(blockDuration) * time.Second

	for i := 0; i < tokenRateLimit; i++ {
		allowed, err := rl.Allow("", "abc123")
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := rl.Allow("", "abc123")
	assert.NoError(t, err)
	assert.False(t, allowed)
}
