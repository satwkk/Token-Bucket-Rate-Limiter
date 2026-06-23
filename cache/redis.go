package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var tryConsumeScript = redis.NewScript(`
local key = KEYS[1]
local maxTokens = tonumber(ARGV[1])
local refillRate = tonumber(ARGV[2])
local currentTime = tonumber(ARGV[3])

local data = redis.call("HMGET", key, "tokens", "last_seen")
local tokens = tonumber(data[1])
local last_seen = tonumber(data[2])

if not tokens or not last_seen then
    tokens = maxTokens
    last_seen = currentTime
end

local tokens_before = tokens

local elapsed = currentTime - last_seen
if elapsed < 0 then elapsed = 0 end

tokens = tokens + (elapsed * refillRate)
if tokens > maxTokens then
    tokens = maxTokens
end

local allowed = 0
if tokens > 1 then
    tokens = tokens - 1
    allowed = 1
end

redis.call("HSET", key, "tokens", tokens, "last_seen", currentTime)
redis.call("EXPIRE", key, 3600)

return {allowed, tokens_before, tokens}
`)

type RedisCache struct {
	client *redis.Client
	config CacheConfig
}

func NewRedisCache(config CacheConfig) Cache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}),
		config: config,
	}
}

func redisFloat(v interface{}) (float64, error) {
	switch n := v.(type) {
	case int64:
		return float64(n), nil
	case float64:
		return n, nil
	case string:
		return strconv.ParseFloat(n, 64)
	default:
		return 0, fmt.Errorf("unexpected redis number type %T", v)
	}
}

func (c *RedisCache) TryConsume(key string) (bool, float64, float64) {
	res, err := tryConsumeScript.Run(
		context.Background(),
		c.client,
		[]string{key},
		c.config.MaxToken,
		c.config.RefillRate,
		float64(time.Now().UnixNano())/1e9,
	).Result()

	if err != nil {
		return false, -1, -1
	}

	vals, ok := res.([]interface{})
	if !ok || len(vals) != 3 {
		return false, -1, -1
	}

	allowedVal, err := redisFloat(vals[0])
	if err != nil {
		return false, -1, -1
	}
	tokensBefore, err := redisFloat(vals[1])
	if err != nil {
		return false, -1, -1
	}
	tokensAfter, err := redisFloat(vals[2])
	if err != nil {
		return false, -1, -1
	}

	return allowedVal == 1, tokensBefore, tokensAfter
}
