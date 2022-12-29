package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellar-file/internal/util/echttputil"
	"github.com/nmluci/stellar-file/pkg/dto"
)

type DoujinBookmarkHandler func(context.Context, *dto.BookQueryDTO) (err error)

func HandleDoujinBookmark(handler DoujinBookmarkHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &dto.BookQueryDTO{}
		if err := c.Bind(req); err != nil {
			fmt.Println(err)
			return echttputil.WriteErrorResponse(c, err)
		}

		if req.StrID != "" {
			req.ID, _ = strconv.ParseInt(req.StrID, 10, 64)
		}

		err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}
