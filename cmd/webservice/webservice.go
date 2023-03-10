package webservice

import (
	"net"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/gostellar"
	"github.com/nmluci/gostellar/pkg/rpc/fileop"
	"github.com/nmluci/stellar-file/cmd/webservice/router"
	inRPC "github.com/nmluci/stellar-file/cmd/webservice/rpc"
	"github.com/nmluci/stellar-file/internal/component"
	"github.com/nmluci/stellar-file/internal/config"
	"github.com/nmluci/stellar-file/internal/interceptor"
	"github.com/nmluci/stellar-file/internal/pubsub"
	"github.com/nmluci/stellar-file/internal/repository"
	"github.com/nmluci/stellar-file/internal/service"
	"github.com/nmluci/stellar-file/internal/worker"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const logTagStartWebservice = "[StartWebservice]"

func Start(conf *config.Config, logger *logrus.Entry) {
	db, err := component.InitMariaDB(&component.InitMariaDBParams{
		Conf:   &conf.MariaDBConfig,
		Logger: logger,
	})

	if err != nil {
		logger.Fatalf("%s initializing maria db: %+v", logTagStartWebservice, err)
	}

	// mongo, err := component.InitMongoDB(&component.InitMongoDBParams{
	// 	Conf:   &conf.MongoDBConfig,
	// 	Logger: logger,
	// })

	// if err != nil {
	// 	logger.Fatalf("%s initializing maria db: %+v", logTagStartWebservice, err)
	// }

	redis, err := component.InitRedis(&component.InitRedisParams{
		Conf:   &conf.RedisConfig,
		Logger: logger,
	})

	if err != nil {
		logger.Fatalf("%s initalizing redis: %+v", logTagStartWebservice, err)
	}

	gs := gostellar.NewGoStellar(gostellar.NewGoStellarParams{
		Logger:      logger,
		ServiceName: conf.ServiceID,
	})

	ec := echo.New()
	ec.HideBanner = true
	ec.HidePort = true

	repo := repository.NewRepository(&repository.NewRepositoryParams{
		Logger:  logger,
		MariaDB: db,
		// MongoDB:    mongo,
		Redis: redis,
	})

	swork := worker.NewWorkerManager(worker.NewWorkerManagerParams{
		Logger:     logger,
		Config:     &conf.WorkerConfig,
		GoStellar:  &gs,
		Repository: repo,
	})

	service := service.NewService(&service.NewServiceParams{
		Logger:     logger,
		Repository: repo,
		StellarRPC: gs.StellarRPC,
		FileWorker: swork,
	})

	psWorker := pubsub.NewFileSub(pubsub.NewFilePubSubParams{
		Logger:  logger,
		Redis:   redis,
		Service: service,
	})

	router.Init(&router.InitRouterParams{
		Logger:  logger,
		Service: service,
		Ec:      ec,
		Conf:    conf,
	})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		swork.StartWorker(5)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Infof("%s starting service, listening to: %s", logTagStartWebservice, conf.ServiceAddress)

		if err := ec.Start(conf.ServiceAddress); err != nil {
			logger.Errorf("%s starting service, cause: %+v", logTagStartWebservice, err)
		}
	}()

	rpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.WithServerInteceptor(service)))
	rpcService := inRPC.Init(service)
	fileop.RegisterStellarFileServer(rpcServer, rpcService)
	reflection.Register(rpcServer)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if conn, err := net.Listen("tcp", conf.RPCAddress); err == nil {
			logger.Infof("%s starting rpc, listening to: %s", logTagStartWebservice, conf.RPCAddress)
			if err := rpcServer.Serve(conn); err != nil {
				logger.Errorf("%s starting rpc, cause: %+v", logTagStartWebservice, err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		psWorker.Listen()
	}()

	wg.Wait()

	swork.StopManager()
}
