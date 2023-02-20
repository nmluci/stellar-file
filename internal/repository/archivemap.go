package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/nmluci/stellar-file/internal/model"
	"github.com/nmluci/stellar-file/pkg/errs"
)

var (
	insertArchivemetaQuery = squirrel.Insert("archivemap").Columns("filename", "filesize", "collection", "created_at")
	getArchivemetaQuery    = squirrel.Select("filename", "filesize", "collection", "created_at").From("archivemap")
)

var (
	tagLoggerFindArchivemetaByID         = "[FindArchivemetaByID]"
	tagLoggerFindArchivemetaByCollection = "[FindArchivemetaByCollection]"
	tagLoggerFindArchivemetaByFilename   = "[FindArchivemetaByFilename]"
	tagLoggerInsertArchivemeta           = "[InsertArchivemeta]"
)

func (r *repository) InsertArchivemeta(ctx context.Context, file *model.ArchiveMap) (err error) {
	query, args, err := insertArchivemetaQuery.Values(file.Filename, file.Filesize, file.Collection, file.CreatedAt).ToSql()
	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerInsertArchivemeta, err)
		return
	}

	_, err = r.mariaDB.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("%s sql err: %+v", tagLoggerInsertArchivemeta, err)
		return
	}

	return
}

func (r *repository) FindArchivemetaByCollection(ctx context.Context, name string) (res []*model.ArchiveMap, err error) {
	query, args, err := getArchivemetaQuery.Where(squirrel.Eq{"collection": name}).ToSql()
	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerFindArchivemetaByCollection, err)
		return
	}

	rows, err := r.mariaDB.QueryxContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("%s sql err: %+v", tagLoggerFindArchivemetaByCollection, err)
		return
	}

	res = []*model.ArchiveMap{}
	for rows.Next() {
		temp := &model.ArchiveMap{}
		err = rows.StructScan(temp)
		if err != nil {
			r.logger.Errorf("%s sql mapping err: %+v", tagLoggerFindArchivemetaByCollection, err)
			return nil, err
		}

		res = append(res, temp)
	}

	return
}

func (r *repository) FindArchivemetaByID(ctx context.Context, id int64) (res *model.ArchiveMap, err error) {
	query, args, err := getArchivemetaQuery.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerFindArchivemetaByID, err)
		return
	}

	err = r.mariaDB.QueryRowxContext(ctx, query, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Errorf("%s sql err: %+v", tagLoggerFindArchivemetaByID, err)
		return
	} else if err == sql.ErrNoRows {
		r.logger.Errorf("%s file not found for id: %d", tagLoggerFindArchivemetaByID, id)
		err = errs.ErrNotFound
		return
	}

	return
}

func (r *repository) FindArchivemetaByFilename(ctx context.Context, filename string) (res *model.ArchiveMap, err error) {
	query, args, err := getArchivemetaQuery.Where(squirrel.Eq{"filename": strings.ReplaceAll(filename, "/", "-")}).ToSql()
	if err != nil {
		r.logger.Errorf("%s squirrel err: %+v", tagLoggerFindArchivemetaByFilename, err)
		return
	}

	res = &model.ArchiveMap{}
	err = r.mariaDB.QueryRowxContext(ctx, query, args...).StructScan(res)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Errorf("%s sql err: %+v", tagLoggerFindArchivemetaByFilename, err)
		return
	} else if err == sql.ErrNoRows {
		r.logger.Errorf("%s file not found for filename: %s", tagLoggerFindArchivemetaByFilename, err)
		err = errs.ErrNotFound
		return
	}

	return
}
