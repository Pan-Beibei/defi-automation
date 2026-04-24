package erc20tokenperiodic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"server/internal/permissions"
	"server/internal/types"
)


const rootAuthority = "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"

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

    if validated.From == nil || *validated.From == "" {
        return nil, fmt.Errorf("from address is required")
    }

    fmt.Printf("收窄之后的权限请求数据: %+v\n", validated)

        // Step 3: 填充默认值（startTime → now if nil）
    populated := PopulatePermission(validated.Permission)

    contracts := permissions.DefaultContracts

      // Step 4: ERC20PeriodTransferEnforcer + ValueLteEnforcer caveats
    localCaveats, err := CreatePermissionCaveats(populated, contracts)
    if err != nil {
        return nil, fmt.Errorf("failed to create permission caveats: %w", err)
    }

    caveats := make([]permissions.Caveat, 0, len(localCaveats)+2)
    for _, c := range localCaveats {
        caveats = append(caveats, permissions.Caveat{
            Enforcer: c.Enforcer,
            Terms:    c.Terms,
            Args:     c.Args,
        })
    }

    // Step 5: expiry rule → TimestampEnforcer caveat
    if expiry, ok := getExpiryTimestamp(validated.Rules); ok {
        caveats = append(caveats, permissions.Caveat{
            Enforcer: contracts.TimestampEnforcer,
            Terms:    timestampTerms(expiry),
            Args:     "0x",
        })
    }

    // Step 6: 获取链上 nonce → NonceEnforcer caveat
    if h.deps.EthRPCURL == "" {
        return nil, fmt.Errorf("EthRPCURL is not configured")
    }

    //getCurrentNonce 使用 eth 库查询
    nonce, err := getCurrentNonce(ctx, h.deps.EthRPCURL, contracts.NonceEnforcer, contracts.DelegationManager, *validated.From)
    if err != nil {
        return nil, fmt.Errorf("failed to get nonce: %w", err)
    }
    caveats = append(caveats, permissions.Caveat{
        Enforcer: contracts.NonceEnforcer,
        Terms:    nonceTerms(nonce),
        Args:     "0x",
    })

		
     // Step 7: 生成随机 salt（crypto/rand，32 bytes）
    saltBuf := make([]byte, 32)
    if _, err := rand.Read(saltBuf); err != nil {
        return nil, fmt.Errorf("failed to generate salt: %w", err)
    }
    salt := "0x" + hex.EncodeToString(saltBuf)

    // Step 8: 组装 PreparedPermission
    return &permissions.PreparedPermission{
        ChainId: h.req.ChainId,
        From:    *validated.From,
        To:      h.req.To,
        UnsignedDelegation: permissions.Delegation{
            Delegate:  h.req.To,
            Delegator: *validated.From,
            Authority: rootAuthority,
            Caveats:   caveats,
            Salt:      salt,
        },
        Caveats:       caveats,
        Justification: populated.Data.Justification,
    }, nil
}

func Register() {
    permissions.Register("erc20-token-periodic", func(
        req types.SubmitPermissionReq,
        deps permissions.HandlerDeps,
    ) permissions.PermissionHandler {
        return &handler{req: req, deps: deps}
    })
}