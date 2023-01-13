package dto

type FileDTO struct {
	ID       int64  `json:"id"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Filesize int64  `json:"filesize"`
}

type FilesDTO struct {
	Data       []*FileDTO `json:"data"`
	Collection string     `json:"collection"`
	Requester  string     `json:"requester"`
}

type FileArchivalDTO struct {
	Filename   string `json:"filename"`
	Collection string `json:"collection"`
	IsFile     bool   `json:"is_file"`
	Requester  string `json:"requester"`
}

type FileQueryDTO struct {
	Collection string `json:"collection"`
	IsBundle   bool   `json:"is_bundle"`
	Requester  string `json:"requester"`
}
