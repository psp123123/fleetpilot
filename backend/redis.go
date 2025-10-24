package backend

import (
	"context"
	"fleetpilot/common/config"
	"fleetpilot/common/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis初始化
func InitRedis() *redis.Client {
	r := config.GlobalCfg.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.Db,
	})

	return rdb

}

// redis写入数据
func RedisSet(k, v string, exp int) error {

	ctx := context.Background()

	// 连接redis
	rdb := InitRedis()
	redisErr := rdb.Set(ctx, k, v, (time.Second * time.Duration(exp))).Err()
	if redisErr != nil {
		logger.Error("Set redis Key error :", redisErr)
		return redisErr
	}

	return nil
}

// redis 获取数据
func RedisGet(k string) (string, error) {
	ctx := context.Background()

	// 连接redis
	rdb := InitRedis()
	val, err := rdb.Get(ctx, k).Result()
	if err != nil {
		logger.Error("Get redis data error: ", err)
		return "", err
	}

	return val, nil
}
