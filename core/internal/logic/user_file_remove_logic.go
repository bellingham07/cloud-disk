package logic

import (
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileRemoveLogic {
	return &UserFileRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileRemoveLogic) UserFileRemove(req *types.UserFileRemoveRequest, userIdentity string) (resp *types.UserFileRemoveReply, err error) {
	// parentId
	parentData := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.ParentIdentity, userIdentity).Get(parentData)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("文件不存在")
	}
	// 修改parent id
	l.svcCtx.Engine.ShowSQL(true)
	sql := "UPDATE user_repository SET parent_id = ? WHERE (identity = ?) AND (`deleted_at` IS NULL)"
	_, err = l.svcCtx.Engine.Exec(sql, parentData.Id, req.Identity)
	return
}
