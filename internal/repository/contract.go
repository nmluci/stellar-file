package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/nmluci/stellar-file/internal/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	InsertBook(ctx context.Context, book *model.Books) (err error)
	FindBook(ctx context.Context, book *model.Books) (res *model.Books, err error)
}

type repository struct {
	mariaDB *sqlx.DB
	mongoDB *mongo.Database
	redis   *redis.Client
	logger  *logrus.Entry
	conf    *repositoryConfig
}

type repositoryConfig struct {
}

type NewRepositoryParams struct {
	Logger  *logrus.Entry
	MariaDB *sqlx.DB
	MongoDB *mongo.Database
	Redis   *redis.Client
}

func NewRepository(params *NewRepositoryParams) Repository {
	return &repository{
		logger:  params.Logger,
		conf:    &repositoryConfig{},
		mariaDB: params.MariaDB,
		mongoDB: params.MongoDB,
		redis:   params.Redis,
	}
}
