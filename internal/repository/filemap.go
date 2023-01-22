package repository

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/nmluci/stellar-file/internal/model"
	"github.com/nmluci/stellar-file/pkg/errs"
)

var (
	insertFileMetaQuery = squirrel.Insert("filemap").Columns("filename", "filesize", "collection", "created_at")
	getFileMetaQuery    = squirrel.Select("filename", "filesize", "collection", "created_at").From("filemap")
)

var (
	tagLoggerFindFilemetaByID         = "[FindFilemetaByID]"
	tagLoggerFindFilemetaByCollection = "[FindFilemetaByCollection]"
	tagLoggerInsertFilemeta           = "[InsertFilemeta]"
)

func (r *repository) InsertFilemeta(ctx context.Context, file *model.FileMap) (err error) {
	query, args, err := insertFileMetaQuery.Values(file.Filename, file.Filesize, file.Collection, file.CreatedAt).ToSql()
	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerInsertFilemeta, err)
		return
	}

	_, err = r.mariaDB.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("%s sql err: %+v", tagLoggerInsertFilemeta, err)
		return
	}

	return
}

func (r *repository) FindFilemetaByCollection(ctx context.Context, name string) (res []*model.FileMap, err error) {
	query, args, err := getFileMetaQuery.Where(squirrel.Eq{"collection": name}).ToSql()
	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerFindFilemetaByCollection, err)
		return
	}

	rows, err := r.mariaDB.QueryxContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("%s sql err: %+v", tagLoggerFindFilemetaByCollection, err)
		return
	}

	res = []*model.FileMap{}
	for rows.Next() {
		temp := &model.FileMap{}
		err = rows.StructScan(temp)
		if err != nil {
			r.logger.Errorf("%s sql mapping err: %+v", tagLoggerFindFilemetaByCollection, err)
			return nil, err
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) FindFilemetaByID(ctx context.Context, id int64) (res *model.FileMap, err error) {
	query, args, err := getFileMetaQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerFindFilemetaByID, err)
		return
	}

	err = r.mariaDB.QueryRowxContext(ctx, query, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Errorf("%s sql err: %+v", tagLoggerFindFilemetaByID, err)
		return
	} else if err == sql.ErrNoRows {
		r.logger.Errorf("%s file not found for id: %d", tagLoggerFindFilemetaByID, id)
		err = errs.ErrNotFound
		return
	}

	return
}
