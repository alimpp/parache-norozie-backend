package util

import (
	"context"
	"crypto/rand"
	"github.com/redis/go-redis/v9"
	"math/big"
	"time"
)

func GenerateRandomString(n int) string {
	const letters = "0123456789"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}

type RedisStore struct {
	Client *redis.Client
	Ctx    context.Context
}

func (r *RedisStore) Get(key string) ([]byte, error) {
	return r.Client.Get(r.Ctx, key).Bytes()
}

// Set saves a session value in Redis
func (r *RedisStore) Set(key string, val []byte, expiration time.Duration) error {
	return r.Client.Set(r.Ctx, key, val, expiration).Err()
}

// Delete removes a session from Redis
func (r *RedisStore) Delete(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

// Reset clears all session data from Redis
func (r *RedisStore) Reset() error {
	return r.Client.FlushDB(r.Ctx).Err()
}

// Close closes the Redis Client connection
func (r *RedisStore) Close() error {
	return r.Client.Close()
}
