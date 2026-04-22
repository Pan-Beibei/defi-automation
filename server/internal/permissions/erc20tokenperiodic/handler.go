package erc20tokenperiodic

import (
	"context"
	"encoding/json"
	"fmt"

	"server/internal/permissions"
	"server/internal/types"
)

// handler 实现 permissions.PermissionHandler 接口
type handler struct {
    req  types.SubmitPermissionReq
    deps permissions.HandlerDeps
}

// Handle 对应 TS: PermissionHandler.handlePermissionRequest() 的无UI简化版
func (h *handler) Handle(ctx context.Context) (*permissions.PreparedPermission, error) {
    // Step 1: 解析并验证 permission.data
    var permData PermissionData
    if err := json.Unmarshal([]byte(h.req.Permission.Data), &permData); err != nil {
        return nil, fmt.Errorf("invalid erc20-token-periodic data: %w", err)
    }
    // ... 验证逻辑
    rules := make([]Rule, 0, len(h.req.Rules))
    for _, r := range h.req.Rules {
        var ruleData map[string]interface{}
        if r.Data != "" {
            if err := json.Unmarshal([]byte(r.Data), &ruleData); err != nil {
                return nil, fmt.Errorf("invalid rule data for %q: %w", r.Type, err)
            }
        }
        rules = append(rules, Rule{
            Type: TypeDescriptor{Name: r.Type},
            Data: ruleData,
        })
    }

    var from *string
    if h.req.From != "" {
        f := h.req.From
        from = &f
    }

    permReq := PermissionRequest{
        ChainID: h.req.ChainId,
        From:    from,
        To:      h.req.To,
        Permission: Permission{
            Type:                TypeDescriptor{Name: h.req.Permission.Type},
            IsAdjustmentAllowed: h.req.IsAdjustmentAllowed,
            Data:                permData,
        },
        Rules: rules,
    }

    // Step 2: 验证 erc20-token-periodic 特有字段 + periodDuration 归一化
    validated, err := ParseAndValidatePermission(permReq)
    if err != nil {
        return nil, err
    }

    fmt.Printf("收窄之后的权限请求数据: %+v\n", validated)

		
    // Step 2: buildContext（调用 token metadata service）
    // Step 3: applyContext（处理 rules）
    // Step 4: createCaveats
    // Step 5: 组装 PreparedPermission

    return nil, nil // 占位，后续 Step 细化
}

// init 在包被导入时自动注册到全局注册表
// 对应 TS: permissionHandlerFactory.ts 中 switch case 的注册
func init() {
    permissions.Register("erc20-token-periodic", func(
        req types.SubmitPermissionReq,
        deps permissions.HandlerDeps,
    ) permissions.PermissionHandler {
        return &handler{req: req, deps: deps}
    })
}