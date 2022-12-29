package handler

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellar-file/internal/util/echttputil"
	"github.com/nmluci/stellar-file/pkg/dto"
)

type RemoteFileHandler func(context.Context, *dto.FileQueryDTO) (err error)

func HandleRemoteFile(handler RemoteFileHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &dto.FileQueryDTO{}
		if err := c.Bind(req); err != nil {
			fmt.Println(err)
			return echttputil.WriteErrorResponse(c, err)
		}

		err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}
