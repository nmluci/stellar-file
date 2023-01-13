package worker

import (
	"github.com/nmluci/stellar-file/internal/config"
	"github.com/sirupsen/logrus"
)

type WorkerManager struct {
	Archival   ArchivalWorker
	Downloader DownloaderWorker
}

type NewWorkerManagerParams struct {
	Logger *logrus.Entry
	Config *config.WorkerConfig
}

func NewWorkerManager(params NewWorkerManagerParams) (wm *WorkerManager) {
	manager := &WorkerManager{
		Archival: NewArchivalWorker(NewArchivalWorkerParams{
			Logger: params.Logger,
			Config: &params.Config.Arc,
		}),
		Downloader: NewDownloaderWorker(NewDownloaderWorkerParams{
			Logger: params.Logger,
			Config: &params.Config.Down,
		}),
	}
	return manager
}

func (wm *WorkerManager) StartWorker(workers int) {
	for i := 0; i < workers; i++ {
		go wm.Archival.Executor(i)
		go wm.Downloader.Executor(i)
	}
}

func (wm *WorkerManager) StopManager() {
	wm.Downloader.StopWorker()
	wm.Archival.StopWorker()
}
