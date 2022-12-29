package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/nmluci/stellar-file/internal/model"
)

var (
	insertBookQuery = squirrel.Insert("tb_books").Columns("book_id", "title", "source", "tags", "created_at", "updated_at")
	selectBookQuery = squirrel.Select("id", "book_id", "title", "source", "tags").From("tb_books")
)

var (
	tagLoggerInsertBook = "[InsertBook]"
	tagLoggerFindBook   = "[FindBook]"
)

func (r *repository) InsertBook(ctx context.Context, book *model.Books) (err error) {
	now := time.Now()
	query, args, err := insertBookQuery.Values(book.BookID, book.Title, book.Source, book.Tags, now, now).ToSql()
	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerInsertBook, err)
		return
	}

	_, err = r.mariaDB.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("%s sql err: %+v", tagLoggerInsertBook, err)
		return
	}

	return
}

func (r *repository) FindBook(ctx context.Context, book *model.Books) (res *model.Books, err error) {
	query, args, err := selectBookQuery.Where(squirrel.And{
		squirrel.Eq{"book_id": book.BookID},
		squirrel.Eq{"source": book.Source},
		squirrel.Eq{"deleted_at": nil},
	}).ToSql()

	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerFindBook, err)
		return
	}

	res = &model.Books{}
	err = r.mariaDB.QueryRowxContext(ctx, query, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Errorf("%s sql err: %+v", tagLoggerFindBook, err)
		return
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}
