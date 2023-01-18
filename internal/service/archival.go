package service

import (
	"context"

	"github.com/nmluci/stellar-file/internal/commonkey"
	"github.com/nmluci/stellar-file/internal/indto"
	"github.com/nmluci/stellar-file/internal/util/ctxutil"
	"github.com/nmluci/stellar-file/internal/util/scopeutil"
	"github.com/nmluci/stellar-file/pkg/dto"
	"github.com/nmluci/stellar-file/pkg/errs"
)

var (
	tagLoggerInsertArchiveJob = "[InsertArchiveJob]"
)

func (s *service) InsertArchiveJob(ctx context.Context, req *dto.FileArchivalDTO) (err error) {
	s.logger.Infoln(ctxutil.GetCtx[indto.UserScopeMap](ctx, commonkey.SCOPE_CTX_KEY))

	if !scopeutil.ValidateScope(ctx, commonkey.FILE_ARCHIVE) && !scopeutil.ValidateScope(ctx, commonkey.FILE_ALL) {
		return errs.ErrNoAccess
	}

	go func() {
		err = s.fileWorker.Archival.InsertJob(req.Collection, req.Filename, req.IsFile)
		if err != nil {
			s.logger.Errorf("%s failed to insert job err: %+v", tagLoggerInsertArchiveJob, err)
		}
	}()

	return
}
