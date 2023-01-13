package service

import (
	"context"

	"github.com/nmluci/stellar-file/pkg/dto"
	"github.com/nmluci/stellar-file/pkg/errs"
)

var (
	tagLoggerInsertDownloadJob = "[InsertDownloadJob]"
)

func (s *service) InsertDownloadJob(ctx context.Context, req *dto.FilesDTO) (err error) {
	if len(req.Data) == 0 {
		s.logger.Errorf("%s data cannot be empty", tagLoggerInsertDownloadJob)
		return errs.ErrBadRequest
	}

	go func() {
		for _, file := range req.Data {
			err = s.fileWorker.Downloader.InsertJob(file.URL, file.Filename, req.Collection)
			if err != nil {
				s.logger.Errorf("%s failed to insert job err: %+v", tagLoggerInsertDownloadJob, err)
				return
			}
		}
	}()

	return
}
