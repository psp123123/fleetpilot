package backend

import (
	"context"
	"fleetpilot/common/config"
	"fleetpilot/common/logger"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisOnce    sync.Once
	redisClient  *redis.Client
	redisInitErr error
)

// 初始化全局 Redis 客户端
func InitRedis() (*redis.Client, error) {
	redisOnce.Do(func() {
		r := config.GlobalCfg.Redis
		redisClient = redis.NewClient(&redis.Options{
			Addr:             r.Address,
			Password:         r.Password,
			DB:               r.Db,
			DisableIndentity: true, // 关键：禁用不支持的功能
			PoolSize:         100,  // 连接池大小
			MinIdleConns:     10,   // 最小空闲连接
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
			redisInitErr = err
			return
		}

		logger.Info("Redis connected successfully")
	})

	return redisClient, redisInitErr
}

// 获取 Redis 客户端
func GetRedis() *redis.Client {
	if redisClient == nil {
		InitRedis()
	}
	return redisClient
}

// Redis 写入数据
func RedisSet(k, v string, exp int) error {
	ctx := context.Background()
	rdb := GetRedis()

	err := rdb.Set(ctx, k, v, time.Second*time.Duration(exp)).Err()
	if err != nil {
		logger.Error("Set redis Key error:", err)
		return err
	}

	return nil
}

// Redis 获取数据
func RedisGet(k string) (string, error) {
	ctx := context.Background()
	rdb := GetRedis()

	val, err := rdb.Get(ctx, k).Result()
	if err != nil && err != redis.Nil { // redis.Nil 是 key 不存在的正常情况
		logger.Error("Get redis data error:", err)
		return "", err
	}

	return val, nil
}
