package service

import (
	"context"

	"github.com/nmluci/gohentai"
	"github.com/nmluci/gostellar"
	"github.com/nmluci/stellar-file/internal/repository"
	"github.com/nmluci/stellar-file/internal/worker"
	"github.com/nmluci/stellar-file/pkg/dto"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)
	AuthenticateSession(ctx context.Context, token string) (access context.Context, err error)
	AuthenticateService(ctx context.Context, name string) (access context.Context, err error)

	InsertDownloadJob(ctx context.Context, req *dto.FilesDTO) (err error)
	InsertArchiveJob(ctx context.Context, req *dto.FileArchivalDTO) (err error)
}

type service struct {
	logger     *logrus.Entry
	conf       *serviceConfig
	repository repository.Repository
	hentailib  gohentai.GoHentai
	stellarRPC *gostellar.StellarRPC
	fileWorker *worker.WorkerManager
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Logger     *logrus.Entry
	Repository repository.Repository
	StellarRPC *gostellar.StellarRPC
	FileWorker *worker.WorkerManager
}

func NewService(params *NewServiceParams) Service {
	return &service{
		logger:     params.Logger,
		conf:       &serviceConfig{},
		repository: params.Repository,
		hentailib:  gohentai.NewHentai(true),
		stellarRPC: params.StellarRPC,
		fileWorker: params.FileWorker,
	}
}
