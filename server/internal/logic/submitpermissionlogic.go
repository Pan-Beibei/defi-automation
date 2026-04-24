// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"

	"server/internal/permissions"
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
	fmt.Printf("权限类型: %s\n", req.Permission.Type)
    
    // 1. 通用字段校验（chainId、to、permission.type、data JSON、expiry）
    if err := permissions.ValidateSubmitPermissionReq(req); err != nil {
        return nil, err
    }

    factory := permissions.NewPermissionHandlerFactory()

    handler, err := factory.Create(req)  // req 是你已验证的 SubmitPermissionReq
    if err != nil {
				fmt.Printf("创建权限处理器出错: %+v\n", err)

        return nil, fmt.Errorf("failed to create permission handler: %w", err)
    }
 

    prepared, err := handler.Handle(l.ctx)
		if err != nil {
			fmt.Printf("处理权限请求出错: %+v\n", err)
				return nil, fmt.Errorf("failed to handle permission request: %w", err)
		}
    

    fmt.Printf("构造好的数据:%+v\n", prepared)

	return &types.SubmitPermissionResp{}, nil
}
