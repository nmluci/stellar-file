package rpcutil

import (
	"context"

	"github.com/nmluci/stellar-file/internal/config"
	"google.golang.org/grpc/metadata"
)

func GenerateMetaContext() context.Context {
	conf := config.Get()
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("st-kagi", conf.StellarConfig.AuthKey))
}

func AppendMetaContext(ctx context.Context) context.Context {
	conf := config.Get()
	return metadata.AppendToOutgoingContext(ctx, "st-kagi", conf.StellarConfig.AuthKey)
}
