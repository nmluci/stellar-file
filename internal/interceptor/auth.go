package interceptor

import (
	"context"
	"time"

	"github.com/nmluci/stellar-file/internal/service"
	"github.com/nmluci/stellar-file/pkg/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func WithServerInteceptor(svc service.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		_ = time.Now()

		ctx, err := authorizeRequest(ctx, svc)
		if err != nil {
			return nil, err
		}

		h, err := handler(ctx, req)

		return h, err
	}
}

func authorizeRequest(ctx context.Context, svc service.Service) (context.Context, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse metadata")
	}

	var err error = nil
	if token, ok := meta["authorization"]; ok {
		ctx, err = svc.AuthenticateSession(ctx, token[0])
	} else if token, ok := meta["st-kagi"]; ok {
		ctx, err = svc.AuthenticateService(ctx, token[0])
	} else {
		return nil, errs.GetErrorRPC(errs.ErrNoAccess)
	}

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to authenticated request")
	}

	return ctx, nil
}
