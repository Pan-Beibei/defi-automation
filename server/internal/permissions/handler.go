package permissions

import (
	"context"
)

// PreparedPermission 是后端处理完成后返回给前端的结果（待签名的 delegation）
// 对应 TS 流程中 #resolveResponse 之前的产物
type PreparedPermission struct {
    ChainId           string     `json:"chainId"`
    From              string     `json:"from"`
    To                string     `json:"to"`
    UnsignedDelegation Delegation `json:"unsignedDelegation"` // 前端用此签名
    Caveats           []Caveat   `json:"caveats"`
    Justification     string     `json:"justification"`
}

// Caveat 对应 delegation-core 中的 Caveat
type Caveat struct {
    Enforcer string `json:"enforcer"`
    Terms    string `json:"terms"` // hex encoded
    Args     string `json:"args"`  // hex encoded, 通常为 "0x"
}

// Delegation 对应待签名的 delegation 结构体
type Delegation struct {
    Delegate  string   `json:"delegate"`
    Delegator string   `json:"delegator"`
    Authority string   `json:"authority"` // ROOT_AUTHORITY = "0xff..."
    Caveats   []Caveat `json:"caveats"`
    Salt      string   `json:"salt"` // hex encoded uint256
}

// HandlerDeps 是所有 handler 共享的依赖
// 对应 TS: PermissionHandlerFactory 中注入的服务
type HandlerDeps struct {
    TokenMetadataService TokenMetadataService // interface, Step后续定义
    EthRPCURL            string // JSON-RPC endpoint, e.g. "https://sepolia.infura.io/v3/..."
}

// DelegationContracts holds enforcer contract addresses (same across all supported chains).
type DelegationContracts struct {
    DelegationManager            string
    NativeTokenStreamingEnforcer string
    ExactCalldataEnforcer        string
    TimestampEnforcer            string
    NonceEnforcer                string
    Erc20PeriodTransferEnforcer  string
    ValueLteEnforcer             string
}

// PermissionHandler 是所有权限处理器必须实现的接口
type PermissionHandler interface {
    Handle(ctx context.Context) (*PreparedPermission, error)
}

// TokenMetadataService 是 token 元数据服务的接口（Step后续实现）
type TokenMetadataService interface {
    GetTokenMetadata(ctx context.Context, chainId int, account string) (*TokenMetadata, error)
}

type TokenMetadata struct {
    Decimals int
    Symbol   string
}