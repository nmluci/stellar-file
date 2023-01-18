package service

import (
	"context"

	"github.com/nmluci/stellar-file/internal/commonkey"
	"github.com/nmluci/stellar-file/internal/config"
	"github.com/nmluci/stellar-file/internal/indto"
	"github.com/nmluci/stellar-file/internal/util/ctxutil"
	"github.com/nmluci/stellar-file/internal/util/rpcutil"
	rpc "github.com/nmluci/stellar-file/pkg/rpc/auth"
)

var (
	tagLoggerAuthenticateSession = "[AuthenticateSession]"
	tagLoggerAuthenticateService = "[AuthenticateService]"
)

func (s *service) AuthenticateSession(ctx context.Context, token string) (access context.Context, err error) {
	conf := config.Get()
	ctx = rpcutil.AppendMetaContext(ctx)

	usr, err := s.stellarRPC.Auth.AuthorizeToken(ctx, &rpc.UserAccess{
		AccessToken: token,
		Requester:   conf.ServiceID,
	})
	if err != nil {
		s.logger.Errorf("%s stellarRPC error: %+v", tagLoggerAuthenticateSession, err)
		return
	}

	scopeMap := indto.UserScopeMap{}
	for _, scope := range usr.UserScope {
		if _, ok := scopeMap[commonkey.FILE_ALL]; !ok && scope == commonkey.FILE_ALL {
			scopeMap[commonkey.FILE_ALL] = true
		}

		if _, ok := scopeMap[commonkey.FILE_ARCHIVE]; !ok && scope == commonkey.FILE_ARCHIVE {
			scopeMap[commonkey.FILE_ARCHIVE] = true
		}

		if _, ok := scopeMap[commonkey.FILE_DOWNLOAD]; !ok && scope == commonkey.FILE_DOWNLOAD {
			scopeMap[commonkey.FILE_DOWNLOAD] = true
		}
	}

	access = ctxutil.WrapCtx(ctx, commonkey.SCOPE_CTX_KEY, scopeMap)
	return
}

func (s *service) AuthenticateService(ctx context.Context, name string) (access context.Context, err error) {
	conf := config.Get()

	ctx = rpcutil.AppendMetaContext(ctx)
	svcMeta, err := s.stellarRPC.Auth.AuthorizeService(ctx, &rpc.ServiceAccess{
		ServiceName: name,
		Requester:   conf.ServiceID,
	})
	if err != nil {
		s.logger.Errorf("%s stellarRPC error: %+v", tagLoggerAuthenticateService, err)
		return
	}

	scopeMap := indto.UserScopeMap{}
	for _, scope := range svcMeta.ServiceScope {
		if _, ok := scopeMap[commonkey.FILE_ALL]; !ok && scope == commonkey.FILE_ALL {
			scopeMap[commonkey.FILE_ALL] = true
		}

		if _, ok := scopeMap[commonkey.FILE_ARCHIVE]; !ok && scope == commonkey.FILE_ARCHIVE {
			scopeMap[commonkey.FILE_ARCHIVE] = true
		}

		if _, ok := scopeMap[commonkey.FILE_DOWNLOAD]; !ok && scope == commonkey.FILE_DOWNLOAD {
			scopeMap[commonkey.FILE_DOWNLOAD] = true
		}
	}

	access = ctxutil.WrapCtx(ctx, commonkey.SCOPE_CTX_KEY, scopeMap)
	return
}
