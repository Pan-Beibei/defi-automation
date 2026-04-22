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

type GetSupportedPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSupportedPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSupportedPermissionsLogic {
	return &GetSupportedPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSupportedPermissionsLogic) GetSupportedPermissions(req *types.GetSupportedPermissionsReq) ( *types.GetSupportedPermissionsResp,  error) {
    // 构建所有链信息列表
    chains := make([]types.ChainInfo, 0, len(svc.NameAndExplorerUrlByChainId))
    for chainId, metadata := range svc.NameAndExplorerUrlByChainId {
        chains = append(chains, types.ChainInfo{
            ChainId: fmt.Sprintf("0x%x", chainId),
            Name:    metadata.Name,
        })
    }

	// 构建权限列表
    permissions := make([]types.PermissionDetail, 0, len(svc.DefaultGatorPermissionToOffer))
    for _, offer := range svc.DefaultGatorPermissionToOffer {
        ruleTypes, ok := svc.SupportedRuleTypes[offer.Type]
        if !ok {
            ruleTypes = []string{}
        }
        permissions = append(permissions, types.PermissionDetail{
            Type:         offer.Type,
            ProposedName: offer.ProposedName,
            Chains:       chains,
            RuleTypes:    ruleTypes,
        })
    }

	return &types.GetSupportedPermissionsResp{
		Permissions: permissions,
	}, nil
}
