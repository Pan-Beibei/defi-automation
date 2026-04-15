import { TextField, Stack } from "@mui/material";
import { toHex, parseUnits } from "viem";
import type { Hex } from "viem";
import { useState, useEffect } from "react";

import type { ERC20TokenStreamPermissionRequest } from "../../types/permissions";
import { BasePermissionFields } from "./BasePermissionFields";

type Props = {
  onChange: (req: ERC20TokenStreamPermissionRequest) => void;
};

const DECIMALS = 6; // USDC

export function ERC20TokenStreamForm({ onChange }: Props) {
  const now = Math.floor(Date.now() / 1000);

  const [tokenAddress, setTokenAddress] = useState<string>(
    "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
  );
  const [initialAmount, setInitialAmount] = useState(
    String(parseUnits("0.5", DECIMALS)),
  );
  const [amountPerSecond, setAmountPerSecond] = useState(
    String(parseUnits("0.5", DECIMALS)),
  );
  const [maxAmount, setMaxAmount] = useState(
    String(parseUnits("2.5", DECIMALS)),
  );
  const [startTime, setStartTime] = useState(String(now));
  const [expiry, setExpiry] = useState(String(now + 60 * 60 * 24 * 30));
  const [justification, setJustification] = useState(
    "This is a very important request for streaming ERC20 token allowance.",
  );
  const [isAdjustmentAllowed, setIsAdjustmentAllowed] = useState(true);

  useEffect(() => {
    onChange({
      type: "erc20-token-stream",
      tokenAddress: tokenAddress as Hex,
      initialAmount: initialAmount.trim() ? toHex(BigInt(initialAmount)) : null,
      amountPerSecond: toHex(BigInt(amountPerSecond || "0")),
      maxAmount: maxAmount.trim() ? toHex(BigInt(maxAmount)) : null,
      startTime: startTime.trim() ? Number(startTime) : null,
      expiry: expiry.trim() ? Number(expiry) : null,
      justification: justification || null,
      isAdjustmentAllowed,
    });
  }, [
    tokenAddress,
    initialAmount,
    amountPerSecond,
    maxAmount,
    startTime,
    expiry,
    justification,
    isAdjustmentAllowed,
    onChange,
  ]);

  return (
    <Stack spacing={2}>
      <TextField
        label="Token Address"
        fullWidth
        value={tokenAddress}
        onChange={(e) => setTokenAddress(e.target.value)}
        placeholder="0x..."
        helperText="ERC20 代币合约地址"
      />
      <Stack direction="row" spacing={2}>
        <TextField
          label={`Initial Amount (最小单位, decimals=${DECIMALS})`}
          fullWidth
          value={initialAmount}
          onChange={(e) => setInitialAmount(e.target.value)}
          helperText="可为空"
        />
        <TextField
          label="Amount Per Second"
          fullWidth
          value={amountPerSecond}
          onChange={(e) => setAmountPerSecond(e.target.value)}
        />
        <TextField
          label="Max Amount"
          fullWidth
          value={maxAmount}
          onChange={(e) => setMaxAmount(e.target.value)}
          helperText="可为空"
        />
      </Stack>
      <BasePermissionFields
        justification={justification}
        startTime={startTime}
        expiry={expiry}
        isAdjustmentAllowed={isAdjustmentAllowed}
        onJustificationChange={setJustification}
        onStartTimeChange={setStartTime}
        onExpiryChange={setExpiry}
        onIsAdjustmentAllowedChange={setIsAdjustmentAllowed}
      />
    </Stack>
  );
}
