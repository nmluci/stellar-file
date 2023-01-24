package router

import (
	"github.com/labstack/echo-contrib/pprof"
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
	params.Ec.GET(PingPath, handler.HandlePing(params.Service.Ping))

	authGroup := params.Ec.Group("")
	authGroup.Use(middleware.AuthorizationMiddleware(params.Service))

	authGroup.POST(DownloadFilePath, handler.HandleDownloadFile(params.Service.InsertDownloadJob))
	authGroup.POST(ArchiveFilePath, handler.HandleArchiveFile(params.Service.InsertArchiveJob))

	pprof.Register(params.Ec)
	// params.Ec.GET(FileIDPath)
	// params.Ec.POST(BookmarkPath, handler.HandleDoujinBookmark(params.Service.BookmarkDoujin))
}
