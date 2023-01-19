package model

type FileMap struct {
	ID         uint64 `db:"id"`
	Filename   string `db:"filename"`
	Filesize   uint64 `db:"filesize"`
	Collection string `db:"collection"`
	CreatedAt  int64  `db:"created_at"`
}

type ArchiveMap struct {
	ID         uint64 `db:"id"`
	Filename   string `db:"filename"`
	Filesize   uint64 `db:"filesize"`
	Collection string `db:"collection"`
	CreatedAt  int64  `db:"created_at"`
}
