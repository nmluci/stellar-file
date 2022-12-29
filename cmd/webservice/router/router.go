package router

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellar-file/cmd/webservice/handler"
	"github.com/nmluci/stellar-file/internal/config"
	"github.com/nmluci/stellar-file/internal/middleware"
	"github.com/nmluci/stellar-file/internal/service"
	"github.com/sirupsen/logrus"
)

type InitRouterParams struct {
	Logger  *logrus.Entry
	Service service.Service
	Ec      *echo.Echo
	Conf    *config.Config
}

func Init(params *InitRouterParams) {
	params.Ec.Use(middleware.AuthorizationMiddleware(params.Service))

	params.Ec.GET(PingPath, handler.HandlePing(params.Service.Ping))
	// params.Ec.GET(FileIDPath)
	// params.Ec.POST(BookmarkPath, handler.HandleDoujinBookmark(params.Service.BookmarkDoujin))
}
