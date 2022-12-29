package component

import (
	"time"

	"github.com/nmluci/stellar-file/internal/config"
	"github.com/nmluci/stellar-file/internal/util/rpcutil"
	rpc "github.com/nmluci/stellar-file/pkg/rpc/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InitStellarRPCParams struct {
	Conf   *config.StellarRPCConfig
	Logger *logrus.Entry
}

type StellarRPCService struct {
	Auth rpc.AuthClient
}

const logTagInitStellarRPC = "[InitStellarRPC]"

func InitStellarRPC(params *InitStellarRPCParams) (srpc *StellarRPCService, err error) {
	srpc = &StellarRPCService{}
	srpc.Auth = initStellarAuth(params)

	params.Logger.Infof("%s stellar-rpc init successfully", logTagInitStellarRPC)
	return
}

func initStellarAuth(params *InitStellarRPCParams) (client rpc.AuthClient) {
	authconn, err := grpc.Dial(params.Conf.AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		params.Logger.Fatalf("%s error init stellar rpc->auth: %+v", logTagInitStellarRPC, err)
		return
	}

	client = rpc.NewAuthClient(authconn)

	conf := config.Get()
	ctx := rpcutil.GenerateMetaContext()

	isConnected := false
	for i := 0; i < 5; i++ {
		_, err := client.Ping(ctx, &rpc.Empty{
			Requester: conf.ServiceID,
		})

		if err == nil {
			isConnected = true
			break
		}
		params.Logger.Errorf("%s error ping rpc->auth: %+v; retrying in 1 second", logTagInitStellarRPC, err)
		time.Sleep(1 * time.Second)
	}

	if !isConnected {
		params.Logger.Fatalf("%s failed to dial stellar-rpc->auth", logTagInitStellarRPC)
	}

	return client
}
