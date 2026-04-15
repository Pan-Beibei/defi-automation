// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"

	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitPermissionLogic {
	return &SubmitPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitPermissionLogic) SubmitPermission(req *types.SubmitPermissionReq) ( *types.SubmitPermissionResp,  error) {
fmt.Printf("SubmitPermissionReq: %+v\n", req)

	return &types.SubmitPermissionResp{}, nil
}
