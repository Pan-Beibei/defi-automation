import { TextField, Stack } from "@mui/material";
import { toHex, parseUnits } from "viem";
import { useState, useEffect, useCallback } from "react";

import type { NativeTokenStreamPermissionRequest } from "../../types/permissions";
import { BasePermissionFields } from "./BasePermissionFields";

type Props = {
  onChange: (req: NativeTokenStreamPermissionRequest) => void;
};

export function NativeTokenStreamForm({ onChange }: Props) {
  const now = Math.floor(Date.now() / 1000);

  const [initialAmount, setInitialAmount] = useState(
    String(parseUnits("0.5", 18)),
  );
  const [amountPerSecond, setAmountPerSecond] = useState(
    String(parseUnits("0.5", 18)),
  );
  const [maxAmount, setMaxAmount] = useState(String(parseUnits("2.5", 18)));
  const [startTime, setStartTime] = useState(String(now));
  const [expiry, setExpiry] = useState(String(now + 60 * 60 * 24 * 30));
  const [justification, setJustification] = useState(
    "This is a very important request for streaming native token allowance.",
  );
  const [isAdjustmentAllowed, setIsAdjustmentAllowed] = useState(true);

  useEffect(() => {
    onChange({
      type: "native-token-stream",
      initialAmount: initialAmount.trim() ? BigInt(initialAmount) : null,
      amountPerSecond: BigInt(amountPerSecond || "0"),
      maxAmount: maxAmount.trim() ? BigInt(maxAmount) : null,
      startTime: startTime.trim() ? Number(startTime) : null,
      expiry: expiry.trim() ? Number(expiry) : null,
      justification: justification || null,
      isAdjustmentAllowed,
    });
  }, [
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
      <Stack direction="row" spacing={2}>
        <TextField
          label="Initial Amount (wei)"
          fullWidth
          value={initialAmount}
          onChange={(e) => setInitialAmount(e.target.value)}
          helperText="可为空，首次可提取的初始金额"
        />
        <TextField
          label="Amount Per Second (wei)"
          fullWidth
          value={amountPerSecond}
          onChange={(e) => setAmountPerSecond(e.target.value)}
          helperText="每秒流速"
        />
        <TextField
          label="Max Amount (wei)"
          fullWidth
          value={maxAmount}
          onChange={(e) => setMaxAmount(e.target.value)}
          helperText="可为空，最大可提取总量"
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
