package worker

import (
	"github.com/nmluci/stellar-file/internal/config"
	"github.com/nmluci/stellar-file/internal/repository"
	"github.com/sirupsen/logrus"
)

type WorkerManager struct {
	Archival   ArchivalWorker
	Downloader DownloaderWorker
	Repository repository.Repository

	TaskChannel chan ArcDownTask
	DoneChannel chan TaskDone
}

type NewWorkerManagerParams struct {
	Logger     *logrus.Entry
	Config     *config.WorkerConfig
	Repository repository.Repository
}

func NewWorkerManager(params NewWorkerManagerParams) (wm *WorkerManager) {
	taskChan := make(chan ArcDownTask, 20)
	doneChan := make(chan TaskDone, 20)

	manager := &WorkerManager{
		TaskChannel: taskChan,
		DoneChannel: doneChan,
		Repository:  params.Repository,
		Archival: NewArchivalWorker(NewArchivalWorkerParams{
			Logger:    params.Logger,
			Config:    &params.Config.Arc,
			Repo:      params.Repository,
			DoneQueue: doneChan,
		}),
		Downloader: NewDownloaderWorker(NewDownloaderWorkerParams{
			Logger:    params.Logger,
			Config:    &params.Config.Down,
			Repo:      params.Repository,
			DoneQueue: doneChan,
		}),
	}
	return manager
}

func (wm *WorkerManager) StartWorker(workers int) {
	go wm.Orchestrator()

	for i := 0; i < workers; i++ {
		go wm.Archival.Executor(i)
		go wm.Downloader.Executor(i)
	}
}

func (wm *WorkerManager) StopManager() {
	wm.Downloader.StopWorker()
	wm.Archival.StopWorker()
}
