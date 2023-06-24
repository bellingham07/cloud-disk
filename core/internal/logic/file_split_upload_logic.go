package logic

import (
	"context"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileSplitUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileSplitUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileSplitUploadLogic {
	return &FileSplitUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileSplitUploadLogic) FileSplitUpload(req *types.FileSplitUploadRequest) (resp *types.FileSplitUploadReply, err error) {
	// todo: add your logic here and delete this line

	return
}
