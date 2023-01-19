package component

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nmluci/stellar-file/internal/config"
	"github.com/sirupsen/logrus"
)

type InitRedisParams struct {
	Conf   *config.RedisConfig
	Logger *logrus.Entry
}

const logTagInitRedis = "[InitRedis]"

func InitRedis(params *InitRedisParams) (db *redis.Client, err error) {
	db = redis.NewClient(&redis.Options{
		Addr:     params.Conf.Address,
		Password: params.Conf.Password,
		DB:       0,
	})

	for i := 20; i > 0; i-- {
		_, err = db.Ping(context.TODO()).Result()
		if err == nil {
			break
		}

		params.Logger.Errorf("%s error init db: %+v, retrying in 1 second", logTagInitRedis, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	params.Logger.Infof("%s redis init succesfully", logTagInitRedis)
	return
}
