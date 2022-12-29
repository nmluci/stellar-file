package service

import (
	"context"
	"strings"

	"github.com/nmluci/gohentai"
	"github.com/nmluci/stellar-file/internal/model"
	"github.com/nmluci/stellar-file/pkg/dto"
	"github.com/nmluci/stellar-file/pkg/errs"
)

var (
	tagLoggerGetDoujinID      = "[GetDoujinByID]"
	tagLoggerGetRandomDoujin  = "[GetRandomDoujin]"
	tagLoggerGetRelatedDoujin = "[GetRelatedDoujin]"
	tagLoggerSearchDoujin     = "[GetDoujinQuery]"
	tagLoggerBookmarkDoujin   = "[BookmarkDoujin]"
)

func (s *service) GetDoujinByNukeID(ctx context.Context, req *dto.BookQueryDTO) (res *dto.BookResponse, err error) {
	book, err := s.hentailib.GetByID(req.ID)
	if err != nil {
		s.logger.Errorf("%s gohentai err: %+v", tagLoggerGetDoujinID, err)
		return
	}

	res = &dto.BookResponse{
		ID:      book.ID,
		MediaID: book.MediaID,
		Title: dto.BookTitle{
			Eng:    book.Title.Eng,
			JP:     book.Title.JP,
			Pretty: book.Title.Pretty,
		},
		Favorites: book.Favorites,
		Thumbnail: book.Thumbnail,
		Cover:     book.Cover,
		Scanlator: book.Scanlator,
		Uploaded:  book.Uploaded,
		EpochTime: book.Epoch,
		Tags:      book.Tags,
		NumPages:  book.NumPages,
	}

	for _, c := range book.Characters {
		res.Characters = append(res.Characters, &dto.BookTag{
			ID:    c.ID,
			Type:  c.Type,
			Name:  c.Name,
			URL:   c.URL,
			Count: c.Count,
		})
	}

	for _, p := range book.Pages {
		res.Pages = append(res.Pages, &dto.BookPage{
			URL:    p.URL,
			Width:  p.Width,
			Height: p.Height,
		})
	}

	for _, rt := range book.RawTags {
		res.RawTags = append(res.RawTags, &dto.BookTag{
			ID:    rt.ID,
			Type:  rt.Type,
			Name:  rt.Name,
			URL:   rt.URL,
			Count: rt.Count,
		})
	}

	return
}

func (s *service) GetRandomDoujin(ctx context.Context) (res *dto.BookResponse, err error) {
	book, err := s.hentailib.Random()
	if err != nil {
		s.logger.Errorf("%s gohentai err: %+v", tagLoggerGetRandomDoujin, err)
		return
	}

	res = &dto.BookResponse{
		ID:      book.ID,
		MediaID: book.MediaID,
		Title: dto.BookTitle{
			Eng:    book.Title.Eng,
			JP:     book.Title.JP,
			Pretty: book.Title.Pretty,
		},
		Favorites: book.Favorites,
		Thumbnail: book.Thumbnail,
		Cover:     book.Cover,
		Scanlator: book.Scanlator,
		Uploaded:  book.Uploaded,
		EpochTime: book.Epoch,
		Tags:      book.Tags,
		NumPages:  book.NumPages,
	}

	for _, c := range book.Characters {
		res.Characters = append(res.Characters, &dto.BookTag{
			ID:    c.ID,
			Type:  c.Type,
			Name:  c.Name,
			URL:   c.URL,
			Count: c.Count,
		})
	}

	for _, p := range book.Pages {
		res.Pages = append(res.Pages, &dto.BookPage{
			URL:    p.URL,
			Width:  p.Width,
			Height: p.Height,
		})
	}

	for _, rt := range book.RawTags {
		res.RawTags = append(res.RawTags, &dto.BookTag{
			ID:    rt.ID,
			Type:  rt.Type,
			Name:  rt.Name,
			URL:   rt.URL,
			Count: rt.Count,
		})
	}

	return
}

func (s *service) GetRelatedDoujin(ctx context.Context, req *dto.BookQueryDTO) (res *dto.BooksResponse, err error) {
	books, err := s.hentailib.Related(req.ID)
	if err != nil {
		s.logger.Errorf("%s gohentai err: %+v", tagLoggerGetRelatedDoujin, err)
		return
	}

	res = &dto.BooksResponse{}
	for _, book := range books {
		temp := &dto.BookResponse{
			ID:      book.ID,
			MediaID: book.MediaID,
			Title: dto.BookTitle{
				Eng:    book.Title.Eng,
				JP:     book.Title.JP,
				Pretty: book.Title.Pretty,
			},
			Favorites: book.Favorites,
			Thumbnail: book.Thumbnail,
			Cover:     book.Cover,
			Scanlator: book.Scanlator,
			Uploaded:  book.Uploaded,
			EpochTime: book.Epoch,
			Tags:      book.Tags,
			NumPages:  book.NumPages,
		}

		for _, c := range book.Characters {
			temp.Characters = append(temp.Characters, &dto.BookTag{
				ID:    c.ID,
				Type:  c.Type,
				Name:  c.Name,
				URL:   c.URL,
				Count: c.Count,
			})
		}

		for _, p := range book.Pages {
			temp.Pages = append(temp.Pages, &dto.BookPage{
				URL:    p.URL,
				Width:  p.Width,
				Height: p.Height,
			})
		}

		for _, rt := range book.RawTags {
			temp.RawTags = append(temp.RawTags, &dto.BookTag{
				ID:    rt.ID,
				Type:  rt.Type,
				Name:  rt.Name,
				URL:   rt.URL,
				Count: rt.Count,
			})
		}

		res.Books = append(res.Books, temp)
	}

	return
}

func (s *service) GetDoujinQuery(ctx context.Context, req *dto.BookQueryDTO) (res *dto.BooksResponse, err error) {
	books, err := s.hentailib.Search(gohentai.SearchParams{
		Query: req.Query,
		Page:  req.Page,
		Sort:  gohentai.Date,
	})
	if err != nil {
		s.logger.Errorf("%s gohentai err: %+v", tagLoggerSearchDoujin, err)
		return
	}

	res = &dto.BooksResponse{}
	for _, book := range books {
		temp := &dto.BookResponse{
			ID:      book.ID,
			MediaID: book.MediaID,
			Title: dto.BookTitle{
				Eng:    book.Title.Eng,
				JP:     book.Title.JP,
				Pretty: book.Title.Pretty,
			},
			Favorites: book.Favorites,
			Thumbnail: book.Thumbnail,
			Cover:     book.Cover,
			Scanlator: book.Scanlator,
			Uploaded:  book.Uploaded,
			EpochTime: book.Epoch,
			Tags:      book.Tags,
			NumPages:  book.NumPages,
		}

		for _, c := range book.Characters {
			temp.Characters = append(temp.Characters, &dto.BookTag{
				ID:    c.ID,
				Type:  c.Type,
				Name:  c.Name,
				URL:   c.URL,
				Count: c.Count,
			})
		}

		for _, p := range book.Pages {
			temp.Pages = append(temp.Pages, &dto.BookPage{
				URL:    p.URL,
				Width:  p.Width,
				Height: p.Height,
			})
		}

		for _, rt := range book.RawTags {
			temp.RawTags = append(temp.RawTags, &dto.BookTag{
				ID:    rt.ID,
				Type:  rt.Type,
				Name:  rt.Name,
				URL:   rt.URL,
				Count: rt.Count,
			})
		}

		res.Books = append(res.Books, temp)
	}

	return
}

func (s *service) BookmarkDoujin(ctx context.Context, req *dto.BookQueryDTO) (err error) {
	if req.ID == 0 {
		s.logger.Errorf("%s ID cannot be empty", tagLoggerBookmarkDoujin)
		return errs.ErrBadRequest
	}

	meta, err := s.hentailib.GetByID(req.ID)
	if err != nil {
		s.logger.Errorf("%s failed to fetch doujin metadata from source err: %+v", tagLoggerBookmarkDoujin, err)
		return
	} else if meta == nil {
		s.logger.Errorf("%s doujin likely already deleted from source", tagLoggerBookmarkDoujin)
		return errs.ErrNotFound
	} else if meta.ID == 0 {
		meta.ID = req.ID
	}

	var exists *model.Books
	if exists, err = s.repository.FindBook(ctx, &model.Books{BookID: meta.ID, Source: "nhentai"}); err != nil {
		s.logger.Errorf("%s failed to check for existing entry for %+v err: %+v", tagLoggerBookmarkDoujin, req.ID, err)
		return
	} else if exists != nil {
		s.logger.Infof("%s book already exists for %+v", tagLoggerBookmarkDoujin, req.ID)
		return errs.ErrDuplicatedResources
	}

	book := &model.Books{
		BookID: meta.ID,
		Title:  meta.Title.Pretty,
		Source: "nhentai",
		Tags:   strings.Join(meta.GetTags(gohentai.Tags), ","),
	}

	err = s.repository.InsertBook(ctx, book)
	if err != nil {
		s.logger.Errorf("%s failed to insert doujin to bookmark collection for %+v: %+v", tagLoggerBookmarkDoujin, book, err)
		return
	}

	return
}
