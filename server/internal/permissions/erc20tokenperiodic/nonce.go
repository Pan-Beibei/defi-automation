package erc20tokenperiodic

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const nonceEnforcerABI = `[{"inputs":[{"name":"delegationManager","type":"address"},{"name":"delegator","type":"address"}],"name":"currentNonce","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

func getCurrentNonce(ctx context.Context, rpcURL, nonceEnforcer, delegationManager, delegator string) (*big.Int, error) {
    client, err := ethclient.DialContext(ctx, rpcURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RPC: %w", err)
    }
    defer client.Close()

    parsedABI, err := abi.JSON(strings.NewReader(nonceEnforcerABI))
    if err != nil {
        return nil, fmt.Errorf("failed to parse NonceEnforcer ABI: %w", err)
    }

    callData, err := parsedABI.Pack("currentNonce",
        common.HexToAddress(delegationManager),
        common.HexToAddress(delegator),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to pack currentNonce args: %w", err)
    }

    to := common.HexToAddress(nonceEnforcer)
    result, err := client.CallContract(ctx, ethereum.CallMsg{
        To:   &to,
        Data: callData,
    }, nil)
    if err != nil {
        return nil, fmt.Errorf("currentNonce call failed: %w", err)
    }

    values, err := parsedABI.Unpack("currentNonce", result)
    if err != nil {
        return nil, fmt.Errorf("failed to unpack currentNonce result: %w", err)
    }
    if len(values) == 0 {
        return nil, fmt.Errorf("currentNonce returned no values")
    }

    nonce, ok := values[0].(*big.Int)
    if !ok {
        return nil, fmt.Errorf("unexpected type from currentNonce: %T", values[0])
    }
    return nonce, nil
}