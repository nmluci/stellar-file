package rpc

import (
	"context"

	"github.com/nmluci/stellar-file/internal/config"
	"github.com/nmluci/stellar-file/pkg/errs"
	"github.com/nmluci/stellar-file/pkg/rpc/fileop"
)

// func (rpc HentaiRPC) Bookmark(ctx context.Context, req *hentai.BookQuery) (res *hentai.Empty, err error) {
// 	if _, ok := config.Get().TrustedService[req.Requester]; !ok {
// 		return nil, errs.GetErrorRPC(errs.ErrNoAccess)
// 	}

// 	if !scopeutil.ValidateScope(ctx, commonkey.NH_SCOPE) {
// 		return nil, errs.GetErrorRPC(errs.ErrNoAccess)
// 	}

// 	params := &dto.BookQueryDTO{
// 		ID: req.GetId(),
// 	}

// 	err = rpc.service.BookmarkDoujin(ctx, params)
// 	if err != nil {
// 		return nil, errs.GetErrorRPC(err)
// 	}

// 	return &hentai.Empty{}, nil
// }

func (rpc FileRPC) Search(ctx context.Context, req *fileop.FileQuery) (res *fileop.Files, err error) {
	if _, ok := config.Get().TrustedService[req.Requester]; !ok {
		return nil, errs.GetErrorRPC(errs.ErrNoAccess)
	}

	return &fileop.Files{}, nil
}

func (rpc FileRPC) Upload(ctx context.Context, req *fileop.Files) (res *fileop.Empty, err error) {
	if _, ok := config.Get().TrustedService[req.Requester]; !ok {
		return nil, errs.GetErrorRPC(errs.ErrNoAccess)
	}

	return &fileop.Empty{}, nil
}
