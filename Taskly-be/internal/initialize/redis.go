package initialize

import (
	"context"
	"fmt"
	"strings"

	"Taskly.com/m/global"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

// InitRedisWithURL khởi tạo Redis client từ chuỗi kết nối (Redis URL)
func InitRedisWithURL(redisURL string) {
    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        global.Logger.Error("Redis URL parse error", zap.Error(err))
        return
    }

    rdb := redis.NewClient(opt)
    global.RedisOpt = asynq.RedisClientOpt{
        Addr:     opt.Addr,
        Password: opt.Password,
        DB:       opt.DB,
    }
    _, err = rdb.Ping(ctx).Result()
    if err != nil {
        global.Logger.Error("Redis initialization Error", zap.Error(err))
        return
    }

    global.Rdb = rdb
}

// Hàm tiện ích để lấy Redis URL từ biến môi trường hoặc config
func InitRedisFromEnvString() {
    redisURL := global.ENVSetting.Redis_Url
    if strings.TrimSpace(redisURL) == "" {
        global.Logger.Error("Redis URL is empty in ENVSetting")
        return
    }
    InitRedisWithURL(redisURL)
}
func InitRedis() {
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", r.Host, r.Port),
		Password: r.Password,
		DB:       r.Database,
		PoolSize: 10,
	})
	global.RedisOpt = asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%v", r.Host, r.Port),
	}
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Error("Redis initialization Error", zap.Error(err))
	}

	global.Rdb = rdb
}
