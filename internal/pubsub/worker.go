package pubsub

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/nmluci/stellar-file/internal/commonkey"
	"github.com/nmluci/stellar-file/internal/indto"
	"github.com/nmluci/stellar-file/internal/service"
	"github.com/nmluci/stellar-file/internal/util/ctxutil"
	"github.com/nmluci/stellar-file/pkg/dto"
	"github.com/sirupsen/logrus"
)

var (
	tagLoggerPBListen = "[PubSub-Listen]"
)

type FilePubSub struct {
	logger  *logrus.Entry
	redis   *redis.Client
	service service.Service
}

type NewFilePubSubParams struct {
	Logger  *logrus.Entry
	Redis   *redis.Client
	Service service.Service
}

func NewFileSub(params NewFilePubSubParams) *FilePubSub {
	return &FilePubSub{
		logger:  params.Logger,
		redis:   params.Redis,
		service: params.Service,
	}
}

func (pb *FilePubSub) Listen() {
	ctx := context.Background()
	ctx = ctxutil.WrapCtx(ctx, commonkey.SCOPE_CTX_KEY, indto.UserScopeMap{
		commonkey.FILE_ALL:      true,
		commonkey.FILE_DOWNLOAD: true,
		commonkey.FILE_ARCHIVE:  true,
	})

	subscriber := pb.redis.Subscribe(ctx, commonkey.TOPICS_FILE_ARC, commonkey.TOPICS_FILE_DOWN)

	defer subscriber.Close()
	for msg := range subscriber.Channel() {
		switch msg.Channel {
		case commonkey.TOPICS_FILE_DOWN:
			var payload dto.FilesDTO
			err := json.Unmarshal([]byte(msg.Payload), &payload)
			if err != nil {
				pb.logger.Errorf("%s payload unmarshaling err: %+v", tagLoggerPBListen, err)
				continue
			}

			err = pb.service.InsertDownloadJob(ctx, &payload)
			if err != nil {
				pb.logger.Errorf("%s batch queue err: %+v", tagLoggerPBListen, err)
				continue
			}

		case commonkey.TOPICS_FILE_ARC:
			var payload dto.FileArchivalDTO
			err := json.Unmarshal([]byte(msg.Payload), &payload)
			if err != nil {
				pb.logger.Errorf("%s payload unmarshaling err: %+v", tagLoggerPBListen, err)
				continue
			}

			err = pb.service.InsertArchiveJob(ctx, &payload)
			if err != nil {
				pb.logger.Errorf("%s batch queue err: %+v", tagLoggerPBListen, err)
				continue
			}

		}
	}
}
