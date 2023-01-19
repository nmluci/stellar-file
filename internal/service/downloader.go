package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nmluci/stellar-file/internal/commonkey"
	"github.com/nmluci/stellar-file/internal/util/scopeutil"
	"github.com/nmluci/stellar-file/pkg/dto"
	"github.com/nmluci/stellar-file/pkg/errs"
)

var (
	tagLoggerInsertDownloadJob = "[InsertDownloadJob]"
)

func (s *service) InsertDownloadJob(ctx context.Context, req *dto.FilesDTO) (err error) {
	if !scopeutil.ValidateScope(ctx, commonkey.FILE_DOWNLOAD) && !scopeutil.ValidateScope(ctx, commonkey.FILE_ALL) {
		return errs.ErrNoAccess
	}

	if len(req.Data) == 0 {
		s.logger.Errorf("%s data cannot be empty", tagLoggerInsertDownloadJob)
		return errs.ErrBadRequest
	}

	go func() {
		taskUUID := uuid.NewString()

		if dataLen := len(req.Data); dataLen > 1 {
			s.fileWorker.DownloadAndArchive(taskUUID, req.Collection, int64(dataLen))
		}

		for _, file := range req.Data {
			err = s.fileWorker.Downloader.InsertJob(taskUUID, file.URL, file.Filename, req.Collection)
			if err != nil {
				s.logger.Errorf("%s failed to insert job err: %+v", tagLoggerInsertDownloadJob, err)
				return
			}
		}
	}()

	return
}
