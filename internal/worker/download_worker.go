package worker

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nmluci/stellar-file/internal/config"
	"github.com/nmluci/stellar-file/internal/repository"
	"github.com/sirupsen/logrus"
)

type DownloaderWorker interface {
	InsertJob(uuid string, target string, filename string, collection string) (err error)
	Executor(id int)
	StopWorker()
}

type downloaderJob struct {
	uuid       string
	target     string
	collection string
	filename   string
}

type downloaderWorker struct {
	client    *http.Client
	wg        *sync.WaitGroup
	logger    *logrus.Entry
	config    *config.DownloaderWorkerConfig
	repo      repository.Repository
	jobQueue  chan downloaderJob
	doneQueue chan TaskDone
}

type NewDownloaderWorkerParams struct {
	Logger    *logrus.Entry
	Repo      repository.Repository
	Config    *config.DownloaderWorkerConfig
	DoneQueue chan TaskDone
}

var (
	tagLoggerDWExecutor   = "[DownloaderWorker-Executor]"
	tagLoggerDWInsertJob  = "[DownloaderWorker-InsertJob]"
	tagLoggerDWStopWorker = "[DownloaderWorker-StopWorker]"
)

func NewDownloaderWorker(params NewDownloaderWorkerParams) (dw DownloaderWorker) {
	httpClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	dw = &downloaderWorker{
		wg:        &sync.WaitGroup{},
		logger:    params.Logger,
		jobQueue:  make(chan downloaderJob, 25),
		client:    httpClient,
		config:    params.Config,
		repo:      params.Repo,
		doneQueue: params.DoneQueue,
	}

	return
}

func (dw *downloaderWorker) InsertJob(uuid string, target string, filename string, collection string) (err error) {
	if target == "" {
		dw.logger.Errorf("%s target cannot be empty", tagLoggerDWInsertJob)
	}

	if filename == "" {
		splitTarget := strings.Split(target, "/")
		filename = splitTarget[len(splitTarget)-1]
	}

	if collection == "" {
		collection = filepath.Join("unknown")
	}

	dw.jobQueue <- downloaderJob{
		uuid:       uuid,
		target:     target,
		filename:   filename,
		collection: filepath.FromSlash(collection),
	}

	return
}

func (dw *downloaderWorker) Executor(id int) {
	dw.logger.Infof("%s Initialized DownloaderWorker id: %d", tagLoggerDWExecutor, id)

	for job := range dw.jobQueue {
		dw.wg.Add(1)
		parentPath := filepath.Join(dw.config.DefaultDir, job.collection)
		fPath := filepath.Join(dw.config.DefaultDir, job.collection, job.filename)

		if fStat, err := os.Stat(fPath); !os.IsNotExist(err) {
			head, err := dw.client.Head(job.target)
			if err != nil {
				dw.logger.Errorf("%s failed to retrieve metadata err: %+v", tagLoggerDWExecutor, err)
				dw.wg.Done()
				dw.doneQueue <- TaskDone{
					UUID:   job.uuid,
					TaskID: TaskDownload,
				}
				continue
			}

			hSize, _ := strconv.ParseInt(head.Header.Get("content-length"), 10, 64)
			if hSize == fStat.Size() {
				dw.wg.Done()
				dw.doneQueue <- TaskDone{
					UUID:   job.uuid,
					TaskID: TaskDownload,
				}
				continue
			}
		}

		req, err := http.NewRequest(http.MethodGet, job.target, nil)
		if err != nil {
			dw.logger.Errorf("%s failed to make new http request err: %+v", tagLoggerDWExecutor, err)
			dw.wg.Done()
			dw.doneQueue <- TaskDone{
				UUID:   job.uuid,
				TaskID: TaskDownload,
			}
			continue
		}

		req.Header.Add("User-Agent", "Stellar-File v1")

		res, err := dw.client.Do(req)
		if err != nil {
			dw.logger.Errorf("%s failed to fetch remote file err: %+v", tagLoggerDWExecutor, err)
			dw.wg.Done()
			dw.doneQueue <- TaskDone{
				UUID:   job.uuid,
				TaskID: TaskDownload,
			}
			continue
		}

		if err := os.MkdirAll(parentPath, os.ModeDir); err != nil && !os.IsExist(err) {
			dw.logger.Errorf("%s failed to make new collection err: %+v", tagLoggerDWExecutor, err)
			dw.wg.Done()
			dw.doneQueue <- TaskDone{
				UUID:   job.uuid,
				TaskID: TaskDownload,
			}
			continue
		}

		fWriter, err := os.OpenFile(fPath, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			dw.logger.Errorf("%s failed to create target file err: %+v", tagLoggerDWExecutor, err)
			dw.wg.Done()
			dw.doneQueue <- TaskDone{
				UUID:   job.uuid,
				TaskID: TaskDownload,
			}
			continue
		}

		_, err = io.Copy(fWriter, res.Body)
		if err != nil {
			fmt.Printf("%s failed to write file err: %+v", tagLoggerDWExecutor, err)
			dw.wg.Done()
			dw.doneQueue <- TaskDone{
				UUID:   job.uuid,
				TaskID: TaskDownload,
			}
			continue
		}

		// if err = dw.repo.InsertFilemeta(context.Background(), &model.FileMap{
		// 	Filename:   job.filename,
		// 	Filesize:   uint64(bWritten),
		// 	Collection: job.collection,
		// 	CreatedAt:  time.Now().UnixMilli(),
		// }); err != nil {
		// 	dw.logger.Errorf("%s failed to log archive into DB err: %+v", tagLoggerDWExecutor, err)
		// 	dw.wg.Done()
		// 	dw.doneQueue <- TaskDone{
		// 		UUID:   job.uuid,
		// 		TaskID: TaskDownload,
		// 	}
		// 	continue
		// }

		fWriter.Close()
		res.Body.Close()

		dw.wg.Done()
		dw.doneQueue <- TaskDone{
			UUID:   job.uuid,
			TaskID: TaskDownload,
		}
	}
}

func (dw *downloaderWorker) StopWorker() {
	dw.wg.Wait()
	dw.logger.Errorf("%s gracefully shutting down worker", tagLoggerDWStopWorker)
	close(dw.jobQueue)
}
