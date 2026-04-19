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

    

	//     switch req.Permission.Type {

    // case "erc20-token-periodic":
    //     // 第一层：结构校验（400 级错误）
    //     data, err := permissions.ParseErc20TokenPeriodicData(req.Permission.Data)
    //     if err != nil {
    //         return nil, err
    //     }
    //     if err := permissions.ValidateErc20TokenPeriodicRequest(req, data); err != nil {
    //         return nil, err
    //     }

    //     // 第二层：字段级校验（前端显示用，需先拿到 token metadata 的 decimals）
    //     validationErrs := permissions.DeriveErc20TokenPeriodicValidationErrors(data, req.Expiry)
    //     if validationErrs.HasErrors() {
    //         // 根据你的错误响应格式返回
    //         return nil, fmt.Errorf("validation failed: %+v", validationErrs)
    //     }

    //     // ... 后续业务逻辑（buildContext、createCaveats 等）

    // default:
    //     return nil, fmt.Errorf("unsupported permission type: %s", req.Permission.Type)
    // }

	return &types.SubmitPermissionResp{}, nil
}
