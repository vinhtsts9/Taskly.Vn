package global

import (
	"database/sql"

	"Taskly.com/m/package/cloudinary"
	"Taskly.com/m/package/kafka"
	"Taskly.com/m/package/setting"
	"github.com/casbin/casbin/v2"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/hibiken/asynq"
	"github.com/pressly/goose/v3/database"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	Config            setting.Config
	RedisOpt          asynq.RedisClientOpt
	ENVSetting 		  setting.ENV
	Logger            *zap.Logger
	PostgreSQL        *sql.DB
	Rdb               *redis.Client
	KafkaProducer     *kafka.Producer
	KafkaConsumer     *kafka.Consumer
	Casbin            *casbin.Enforcer
	Cloudinary        *cloudinary.CloudinaryService
	Store             database.Store
	Elasticsearch     *elasticsearch.Client
)
