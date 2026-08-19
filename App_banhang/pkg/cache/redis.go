package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheService struct {
	ctx context.Context
	rdb *redis.Client
}

func NewRedisCacheService(rdb *redis.Client) *RedisCacheService {
	return &RedisCacheService{
		ctx: context.Background(),
		rdb: rdb,
	}
}

func (r *RedisCacheService) Get(key string, dest any) error {
	data, err := r.rdb.Get(r.ctx, key).Result()
	if err == redis.Nil || err != nil {
		return nil
	}
	return json.Unmarshal([]byte(data), dest)
}

func (r *RedisCacheService) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.rdb.Set(r.ctx, key, data, ttl).Err()
}

func (r *RedisCacheService) Clear(pattern string) error {
	cusor := uint64(0)
	for {
		keys, nextCusor, err := r.rdb.Scan(r.ctx, cusor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			r.rdb.Del(r.ctx, keys...)
		}
		cusor = nextCusor
		if cusor == 0 {
			break
		}
	}
	return nil
}
