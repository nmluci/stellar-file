package service

import (
	"context"

	"github.com/nmluci/gohentai"
	"github.com/nmluci/stellar-file/internal/component"
	"github.com/nmluci/stellar-file/internal/repository"
	"github.com/nmluci/stellar-file/pkg/dto"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)
	GetDoujinByNukeID(ctx context.Context, req *dto.BookQueryDTO) (res *dto.BookResponse, err error)
	GetRandomDoujin(ctx context.Context) (res *dto.BookResponse, err error)
	GetRelatedDoujin(ctx context.Context, req *dto.BookQueryDTO) (res *dto.BooksResponse, err error)
	GetDoujinQuery(ctx context.Context, req *dto.BookQueryDTO) (res *dto.BooksResponse, err error)

	BookmarkDoujin(ctx context.Context, req *dto.BookQueryDTO) (err error)

	AuthenticateSession(ctx context.Context, token string) (access context.Context, err error)
	AuthenticateService(ctx context.Context, name string) (access context.Context, err error)
}

type service struct {
	logger     *logrus.Entry
	conf       *serviceConfig
	repository repository.Repository
	hentailib  gohentai.GoHentai
	stellarRPC *component.StellarRPCService
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Logger     *logrus.Entry
	Repository repository.Repository
	StellarRPC *component.StellarRPCService
}

func NewService(params *NewServiceParams) Service {
	return &service{
		logger:     params.Logger,
		conf:       &serviceConfig{},
		repository: params.Repository,
		hentailib:  gohentai.NewHentai(true),
		stellarRPC: params.StellarRPC,
	}
}
