package worker

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/nmluci/gostellar"
	"github.com/nmluci/gostellar/pkg/dto"
	"github.com/nmluci/stellar-file/internal/repository"
	"github.com/sirupsen/logrus"
)

type NotificationWorker interface {
	InsertJob(uuid string, collection string, filename string, jobType int64, sum int64) (err error)
	Executor()
	StopWorker()
}

type notificationJob struct {
	uuid       string
	collection string
	filename   string
	jobType    int64
	sum        int64
}

type notificationWorker struct {
	wg        *sync.WaitGroup
	logger    *logrus.Entry
	repo      repository.Repository
	gostellar *gostellar.GoStellar

	jobQueue chan notificationJob
	lastPush time.Time
}

type NewNotificationWorkerParams struct {
	Logger    *logrus.Entry
	Repo      repository.Repository
	GoStellar *gostellar.GoStellar
}

var (
	tagLoggerNWExecutor   = "[NotificationWorker-Executor]"
	tagLoggerNWStopWorker = "[NotificationWorker-StopWorker]"
)

func NewNotificationWorker(params NewNotificationWorkerParams) (nw NotificationWorker) {
	nw = &notificationWorker{
		wg:        &sync.WaitGroup{},
		logger:    params.Logger,
		repo:      params.Repo,
		jobQueue:  make(chan notificationJob, 10),
		gostellar: params.GoStellar,
		lastPush:  time.Now(),
	}

	return
}

func (nw *notificationWorker) InsertJob(uuid string, collection string, filename string, jobType int64, sum int64) (err error) {
	nw.jobQueue <- notificationJob{
		uuid:       uuid,
		collection: collection,
		filename:   filename,
		jobType:    jobType,
		sum:        sum,
	}

	return
}

func (nw *notificationWorker) Executor() {
	nw.logger.Infof("%s Initialized NotificatinWorker", tagLoggerNWExecutor)

	for {
		if len(nw.jobQueue) < 10 || !nw.lastPush.Add(5*time.Minute).After(time.Now()) {
			continue
		}

		baseWebhookMeta := &dto.DiscordWebhookMeta{
			Username: "Natsumi-chan",
		}

		qLen := len(nw.jobQueue)
		if qLen >= 3 {
			baseWebhookMeta.Embeds = append(baseWebhookMeta.Embeds, dto.DiscoedEmbeds{
				Title:       "Stellar-File Download and Archive Summary",
				Description: fmt.Sprintf("summary for last %d jobs", qLen),
				Color:       "13421823",
				Footer: dto.DiscordFooter{
					Text: "Stellar-File Download-Archival Worker Summary | Stellar-MS",
				},
				Author: dto.DiscordAuther{
					Name: "Stellar-File by Natsumi-chan",
				},
			})

			fields := []dto.DiscordField{}
			for i := 0; i < qLen; i++ {
				data := <-nw.jobQueue
				temp := dto.DiscordField{
					Value: data.filename,
				}

				if data.jobType == TaskArchive {
					temp.Name = "archival"
				} else if data.jobType == TaskDownload {
					temp.Name = "download"
				}

				fields = append(fields, temp)
			}

			baseWebhookMeta.Embeds[0].Fields = fields
		} else {
			for i := 0; i < qLen; i++ {
				data := <-nw.jobQueue
				temp := dto.DiscoedEmbeds{
					Title:       "Stellar-File Download and Archive Summary",
					Description: fmt.Sprintf("uuid: %s", data.uuid),
					Color:       "13421823",
					Footer: dto.DiscordFooter{
						Text: "Stellar-File Download-Archival Worker Summary | Stellar-MS",
					},
					Author: dto.DiscordAuther{
						Name: "Stellar-File by Natsumi-chan",
					},
				}

				if data.jobType == TaskArchive {
					temp.Fields = append(temp.Fields, dto.DiscordField{
						Name:  "Jobs",
						Value: "archival",
					})

					fMeta, err := nw.repo.FindArchivemetaByFilename(context.Background(), data.filename)
					if err != nil {
						nw.logger.Warnf("%s failed to retrieve archive info from DB err: %+v", tagLoggerNWExecutor, err)
					} else {
						temp.Fields = append(temp.Fields, dto.DiscordField{
							Name:  "Archive sizes",
							Value: fmt.Sprintf("%s MB", strconv.FormatFloat(float64(fMeta.Filesize/uint64(math.Pow(2, 20))), 'f', 2, 64)),
						})
					}

				} else if data.jobType == TaskDownload {
					temp.Fields = append(temp.Fields, dto.DiscordField{
						Name:  "Jobs",
						Value: "download",
					})

					temp.Fields = append(temp.Fields, dto.DiscordField{
						Name:  "Files",
						Value: fmt.Sprintf("%d files", data.sum),
					})
				}

				temp.Fields = append(temp.Fields, dto.DiscordField{
					Name:  "Collection",
					Value: data.collection,
				})

				baseWebhookMeta.Embeds = append(baseWebhookMeta.Embeds, temp)
			}
		}

		if err := nw.gostellar.Notification.Discord.Notify(baseWebhookMeta); err != nil {
			nw.logger.Warnf("%s failed to send notification err: %+v", tagLoggerNWExecutor, err)
		}
		nw.lastPush = time.Now()
	}

}

func (nw *notificationWorker) StopWorker() {
	nw.wg.Wait()
	nw.logger.Errorf("%s gracefully shutting notification worker", tagLoggerNWStopWorker)
	close(nw.jobQueue)
}
