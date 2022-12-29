package model

type Books struct {
	ID     int64  `db:"id"`
	BookID int64  `db:"book_id"`
	Title  string `db:"title"`
	Source string `db:"source"`
	Tags   string `db:"tags"`
}
