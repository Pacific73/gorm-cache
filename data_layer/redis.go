package data_layer

import "github.com/go-redis/redis"

type RedisLayer struct {
	client *redis.Client
}

func (r *RedisLayer) a() {
	r.client.Get("").Err()
}
