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
	InsertFilemeta(ctx context.Context, file *model.FileMap) (err error)
	FindFilemetaByCollection(ctx context.Context, name string) (res []*model.FileMap, err error)
	FindFilemetaByID(ctx context.Context, id int64) (res *model.FileMap, err error)

	InsertArchivemeta(ctx context.Context, file *model.ArchiveMap) (err error)
	FindArchivemetaByCollection(ctx context.Context, name string) (res []*model.ArchiveMap, err error)
	FindArchivemetaByID(ctx context.Context, id int64) (res *model.ArchiveMap, err error)
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
