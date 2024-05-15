package data

import (
	"context"
	"time"
	"userDemo/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 第四步注入
// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGormDB, NewRedis, NewUserRepo)

// 第一步定义Data结构体
type Data struct {
	gormDB *gorm.DB
	redis  *redis.Client
}

// 第 3 步，初始化 Data
func NewData(logger log.Logger, db *gorm.DB, redis *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	return &Data{gormDB: db, redis: redis}, cleanup, nil
}

// 第 2 步初始化 gorm
func NewGormDB(c *conf.Data) (*gorm.DB, error) {
	dsn := c.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(150)
	sqlDB.SetConnMaxLifetime(time.Second * 25)
	return db, err
}

// 初始化redis
func NewRedis(data *conf.Data) *redis.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client := redis.NewClient(&redis.Options{
		Addr: data.Redis.Addr,    // no password set
		DB:   int(data.Redis.Db), // use default DB
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return client
}
