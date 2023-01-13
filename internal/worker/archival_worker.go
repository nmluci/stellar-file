package worker

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/nmluci/stellar-file/internal/config"
	"github.com/nmluci/stellar-file/pkg/errs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ArchivalWorker interface {
	InsertJob(path string, filename string, isFile bool) (err error)
	Executor(id int)
	StopWorker()
}

type archivalJob struct {
	path     string
	filename string
	isFile   bool
}

type archivalWorker struct {
	wg       *sync.WaitGroup
	logger   *logrus.Entry
	jobQueue chan archivalJob
	config   *config.ArchivalWorkerConfig
}

type NewArchivalWorkerParams struct {
	Logger *logrus.Entry
	Config *config.ArchivalWorkerConfig
}

var (
	tagLoggerAWInsertJob  = "[ArchivalWorker-InsertJob]"
	tagLoggerAWExecutor   = "[ArchivalWorker-Executor]"
	tagLoggerAWStopWorker = "[ArchivalWorker-StopWorker]"
)

func NewArchivalWorker(params NewArchivalWorkerParams) (aw ArchivalWorker) {
	aw = &archivalWorker{
		wg:       &sync.WaitGroup{},
		logger:   params.Logger,
		config:   params.Config,
		jobQueue: make(chan archivalJob, 10),
	}

	return
}

func (aw *archivalWorker) InsertJob(path string, filename string, isFile bool) (err error) {
	if path == "" || filename == "" {
		aw.logger.Errorf("%s path and/or filename cannot be empty", tagLoggerAWInsertJob)
		return errs.ErrBadRequest
	}

	aw.jobQueue <- archivalJob{
		path:     path,
		isFile:   isFile,
		filename: filename,
	}

	return
}

func (aw *archivalWorker) Executor(id int) {
	aw.logger.Infof("%s Initialized ArchivalWorker id: %d", tagLoggerAWExecutor, id)

	for job := range aw.jobQueue {
		aw.wg.Add(1)

		conf := config.Get().WorkerConfig
		job.path = filepath.Join(conf.Down.DefaultDir, job.path)

		if err := os.Mkdir(aw.config.DefaultDir, fs.ModeDir); err != nil && !os.IsExist(err) {
			aw.logger.Errorf("%s failed to create new default directory err: %+v", tagLoggerAWExecutor, err)
			continue
		}

		fWriter, err := os.OpenFile(filepath.Join(aw.config.DefaultDir, fmt.Sprintf("%s.tar", job.filename)), os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			aw.logger.Errorf("%s failed to create new tar file err: %+v", tagLoggerAWExecutor, err)
			aw.wg.Done()
			continue
		}

		tarWriter := tar.NewWriter(fWriter)

		if !job.isFile {
			files, metas, err := aw.collectFile(job.path)
			if err != nil {
				aw.logger.Errorf("%s an error occurred while collecting files from path err: %+v", tagLoggerAWExecutor, err)
				aw.wg.Done()
				continue
			}

			for idx, file := range files {
				meta, _ := metas[idx].Info()

				tarWriter.WriteHeader(&tar.Header{
					Name:    meta.Name(),
					Size:    meta.Size(),
					ModTime: meta.ModTime(),
				})

				if _, err = io.Copy(tarWriter, file); err != nil {
					aw.logger.Errorf("%s failed to write file err: %+v", tagLoggerAWExecutor, err)
					aw.wg.Done()
					continue
				}
			}
		} else {
			meta, err := os.Stat(job.path)
			if err != nil {
				aw.logger.Errorf("%s failed to fetch file meta err: %+v", tagLoggerAWExecutor, err)
				aw.wg.Done()
				continue
			}

			file, err := os.OpenFile(job.path, os.O_RDWR, 0666)
			if err != nil {
				aw.logger.Errorf("%s failed to read filepath: %s err: %+v", tagLoggerAWExecutor, job.path, err)
				aw.wg.Done()
				continue
			}

			tarWriter.WriteHeader(&tar.Header{
				Name:    meta.Name(),
				Size:    meta.Size(),
				ModTime: meta.ModTime(),
			})

			if _, err = io.Copy(tarWriter, file); err != nil {
				aw.logger.Errorf("%s failed to write file err: %+v", tagLoggerAWExecutor, err)
				aw.wg.Done()
				continue
			}
		}

		if err = tarWriter.Close(); err != nil {
			aw.logger.Errorf("%s failed to close tar filestream err: %+v", tagLoggerAWExecutor, err)
			aw.wg.Done()
			continue
		}

		if err = fWriter.Close(); err != nil {
			aw.logger.Errorf("%s failed to close filestream err: %+v", tagLoggerAWExecutor, err)
			aw.wg.Done()
			continue
		}

		if err = os.RemoveAll(job.path); err != nil {
			aw.logger.Warnf("%s unable to delete file err: %+v", tagLoggerAWExecutor, err)
			aw.wg.Done()
			continue
		}

		aw.wg.Done()
	}
}

func (aw *archivalWorker) collectFile(filedir string) (files []*os.File, filemeta []fs.DirEntry, err error) {
	filemeta, err = os.ReadDir(filedir)
	if err != nil {
		err = errors.Errorf("failed to read filepath: %s err: %+v", filedir, err)
		return
	}

	for _, file := range filemeta {
		if file.IsDir() {
			err = errors.New("only upto 1 level of nested relative path allowed")
			return nil, nil, err
		}

		data, err := os.OpenFile(filepath.Join(filedir, file.Name()), os.O_RDWR, 0666)
		if err != nil {
			err = errors.Errorf("failed to read file: %s err: %+v", file.Name(), err)
			return nil, nil, err
		}

		files = append(files, data)
	}

	return
}

func (aw *archivalWorker) StopWorker() {
	aw.wg.Wait()
	aw.logger.Errorf("%s gracefully shutting down worker", tagLoggerAWStopWorker)
	close(aw.jobQueue)
}
