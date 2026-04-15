import { TextField, Stack } from "@mui/material";
import { toHex, parseUnits } from "viem";
import type { Hex } from "viem";
import { useState, useEffect } from "react";

import type { ERC20TokenPeriodicPermissionRequest } from "../../types/permissions";
import { BasePermissionFields } from "./BasePermissionFields";

type Props = {
  onChange: (req: ERC20TokenPeriodicPermissionRequest) => void;
};

const DECIMALS = 6;

export function ERC20TokenPeriodicForm({ onChange }: Props) {
  const now = Math.floor(Date.now() / 1000);

  const [tokenAddress, setTokenAddress] = useState<string>(
    "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
  );
  const [periodAmount, setPeriodAmount] = useState(
    String(parseUnits("1", DECIMALS)),
  );
  const [periodDuration, setPeriodDuration] = useState("2592000");
  const [startTime, setStartTime] = useState(String(now));
  const [expiry, setExpiry] = useState(String(now + 60 * 60 * 24 * 30));
  const [justification, setJustification] = useState(
    "This is a very important request for periodic ERC20 token allowance.",
  );
  const [isAdjustmentAllowed, setIsAdjustmentAllowed] = useState(true);

  useEffect(() => {
    onChange({
      type: "erc20-token-periodic",
      tokenAddress: tokenAddress as Hex,
      periodAmount: toHex(BigInt(periodAmount || "0")),
      periodDuration: Number(periodDuration),
      startTime: startTime.trim() ? Number(startTime) : null,
      expiry: expiry.trim() ? Number(expiry) : null,
      justification: justification || null,
      isAdjustmentAllowed,
    });
  }, [
    tokenAddress,
    periodAmount,
    periodDuration,
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
          label={`Period Amount (decimals=${DECIMALS})`}
          fullWidth
          value={periodAmount}
          onChange={(e) => setPeriodAmount(e.target.value)}
        />
        <TextField
          label="Period Duration (seconds)"
          fullWidth
          type="number"
          value={periodDuration}
          onChange={(e) => setPeriodDuration(e.target.value)}
          helperText="2592000 = 30 天"
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
