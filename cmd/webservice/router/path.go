package router

const (
	basePath         = "/v1"
	PingPath         = basePath + "/ping"
	FileIDPath       = basePath + "/:id"
	DownloadFilePath = basePath + "/download"
	ArchiveFilePath  = basePath + "/archive"
)
