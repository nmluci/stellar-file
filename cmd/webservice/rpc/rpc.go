package rpc

import (
	fRPC "github.com/nmluci/gostellar/pkg/rpc/fileop"
	insvc "github.com/nmluci/stellar-file/internal/service"
)

type FileRPC struct {
	fRPC.UnimplementedStellarFileServer
	service insvc.Service
}

func Init(svc insvc.Service) fRPC.StellarFileServer {
	return &FileRPC{
		service: svc,
	}
}
