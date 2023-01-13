package config

type WorkerConfig struct {
	Arc  ArchivalWorkerConfig
	Down DownloaderWorkerConfig
}

type ArchivalWorkerConfig struct {
	DefaultDir string
}

type DownloaderWorkerConfig struct {
	DefaultDir string
}
