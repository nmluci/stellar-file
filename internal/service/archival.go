package service

import (
	"context"

	"github.com/nmluci/stellar-file/pkg/dto"
)

var (
	tagLoggerInsertArchiveJob = "[InsertArchiveJob]"
)

func (s *service) InsertArchiveJob(ctx context.Context, req *dto.FileArchivalDTO) (err error) {
	go func() {
		err = s.fileWorker.Archival.InsertJob(req.Collection, req.Filename, req.IsFile)
		if err != nil {
			s.logger.Errorf("%s failed to insert job err: %+v", tagLoggerInsertArchiveJob, err)
		}
	}()

	return
}
