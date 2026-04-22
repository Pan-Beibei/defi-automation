package erc20tokenperiodic

import (
	"fmt"
	"math/big"
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
    if !addressRegex.MatchString(tokenAddress) {
        return "", fmt.Errorf("invalid tokenAddress: must be a valid address")
    }
    periodAmount := new(big.Int)
    hexStr := strings.TrimPrefix(periodAmountHex, "0x")
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
    amtHex := padLeft64(periodAmount.Text(16))
    durHex := fmt.Sprintf("%064x", uint64(periodDuration))
    dateHex := fmt.Sprintf("%064x", uint64(startDate))

    return "0x" + addrHex + amtHex + durHex + dateHex, nil
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
func CreatePermissionCaveats(perm PopulatedPermission, contracts DelegationContracts) ([]Caveat, error) {
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