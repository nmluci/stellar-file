package rpc

import (
	insvc "github.com/nmluci/stellar-file/internal/service"
	fRPC "github.com/nmluci/stellar-file/pkg/rpc/fileop"
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
