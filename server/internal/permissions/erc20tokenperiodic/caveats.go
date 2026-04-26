package erc20tokenperiodic

import (
	"fmt"
	"math/big"
	"server/internal/permissions"
	"strings"
)

// erc20PeriodTransferTerms encodes the terms for ERC20PeriodTransferEnforcer.
//
// Tightly-packed layout (matches createERC20TokenPeriodTransferTerms in delegation-core):
//
//	tokenAddress  : 20 bytes  (40 hex chars, no padding — raw address bytes)
//	periodAmount  : 32 bytes  (64 hex chars, big-endian uint256)
//	periodDuration: 32 bytes  (64 hex chars, big-endian uint256)
//	startDate     : 32 bytes  (64 hex chars, big-endian uint256)
//
// Total: 116 bytes, returned as a 0x-prefixed hex string.
func erc20PeriodTransferTerms(tokenAddress, periodAmountHex string, periodDuration, startDate int64) (string, error) {
    // validate token address
    if !addressRegex.MatchString(tokenAddress) {
        return "", fmt.Errorf("invalid tokenAddress: must be a valid address")
    }
    periodAmount := new(big.Int)
    hexStr := strings.TrimPrefix(periodAmountHex, "0x")
    // parse input hex
    if _, ok := periodAmount.SetString(hexStr, 16); !ok || periodAmount.Sign() <= 0 {
        return "", fmt.Errorf("invalid periodAmount: must be a positive hex integer")
    }
    if periodDuration <= 0 {
        return "", fmt.Errorf("invalid periodDuration: must be a positive number")
    }
    if startDate <= 0 {
        return "", fmt.Errorf("invalid startDate: must be a positive number")
    }

    addrHex := strings.ToLower(strings.TrimPrefix(tokenAddress, "0x")) // 40 hex chars
    amtHex := padLeft64(periodAmount.Text(16)) // output hex padded to 64 chars (32 bytes)
    durHex := fmt.Sprintf("%064x", uint64(periodDuration))
    dateHex := fmt.Sprintf("%064x", uint64(startDate))

    return "0x" + addrHex + amtHex + durHex + dateHex, nil
}

// timestampTerms 编码 TimestampEnforcer 的 terms：
// 0x + afterThreshold(16B=0) + beforeThreshold(16B=expiry)
// 对应 TS: delegation-core createTimestampTerms()
func timestampTerms(expiry int64) string {
    after := fmt.Sprintf("%032x", 0)
    before := fmt.Sprintf("%032x", uint64(expiry))
    return "0x" + after + before
}

// nonceTerms 编码 NonceEnforcer 的 terms：0x + nonce padded to 32 bytes
// 对应 TS: delegation-core createNonceTerms()
func nonceTerms(nonce *big.Int) string {
    h := nonce.Text(16)
    if len(h) < 64 {
        h = strings.Repeat("0", 64-len(h)) + h
    }
    return "0x" + h
}

// valueLteTerms encodes the terms for ValueLteEnforcer with maxValue = 0.
// Returns 0x + 64 zero hex chars (32-byte zero).
// Matches createValueLteTerms({ maxValue: 0n }) in delegation-core.
func valueLteTerms() string {
    return fmt.Sprintf("0x%064x", 0)
}

// padLeft64 left-pads a hex string (no 0x prefix) to 64 characters.
func padLeft64(h string) string {
    if len(h) >= 64 {
        return h
    }
    return strings.Repeat("0", 64-len(h)) + h
}

// CreatePermissionCaveats builds the delegation caveats for an erc20-token-periodic permission.
//
// Returns two caveats in order:
//  1. ERC20PeriodTransferEnforcer — enforces tokenAddress, periodAmount cap, periodDuration, and startDate.
//  2. ValueLteEnforcer            — caps msg.value at 0 (no native token allowed alongside the delegation).
//
// Matches createPermissionCaveats() in caveats.ts.
func CreatePermissionCaveats(perm PopulatedPermission, contracts permissions.DelegationContracts) ([]Caveat, error) {
    erc20Terms, err := erc20PeriodTransferTerms(
        perm.Data.TokenAddress,
        perm.Data.PeriodAmount,
        perm.Data.PeriodDuration,
        perm.Data.StartTime,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create ERC20 period transfer terms: %w", err)
    }

    return []Caveat{
        {
            Enforcer: contracts.Erc20PeriodTransferEnforcer,
            Terms:    erc20Terms,
            Args:     "0x",
        },
        {
            Enforcer: contracts.ValueLteEnforcer,
            Terms:    valueLteTerms(),
            Args:     "0x",
        },
    }, nil
}