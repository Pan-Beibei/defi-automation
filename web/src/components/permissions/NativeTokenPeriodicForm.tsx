import { TextField, Stack } from "@mui/material";
import { toHex, parseUnits } from "viem";
import { useState, useEffect } from "react";

import type { NativeTokenPeriodicPermissionRequest } from "../../types/permissions";
import { BasePermissionFields } from "./BasePermissionFields";

type Props = {
  onChange: (req: NativeTokenPeriodicPermissionRequest) => void;
};

export function NativeTokenPeriodicForm({ onChange }: Props) {
  const now = Math.floor(Date.now() / 1000);

  const [periodAmount, setPeriodAmount] = useState(String(parseUnits("1", 18)));
  const [periodDuration, setPeriodDuration] = useState("2592000"); // 30 days
  const [startTime, setStartTime] = useState(String(now));
  const [expiry, setExpiry] = useState(String(now + 60 * 60 * 24 * 30));
  const [justification, setJustification] = useState(
    "This is a very important request for periodic native token allowance.",
  );
  const [isAdjustmentAllowed, setIsAdjustmentAllowed] = useState(true);

  useEffect(() => {
    onChange({
      type: "native-token-periodic",
      periodAmount: toHex(BigInt(periodAmount || "0")),
      periodDuration: Number(periodDuration),
      startTime: startTime.trim() ? Number(startTime) : null,
      expiry: expiry.trim() ? Number(expiry) : null,
      justification: justification || null,
      isAdjustmentAllowed,
    });
  }, [
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
      <Stack direction="row" spacing={2}>
        <TextField
          label="Period Amount (wei)"
          fullWidth
          value={periodAmount}
          onChange={(e) => setPeriodAmount(e.target.value)}
          helperText="每个周期可提取的金额"
        />
        <TextField
          label="Period Duration (seconds)"
          fullWidth
          type="number"
          value={periodDuration}
          onChange={(e) => setPeriodDuration(e.target.value)}
          helperText="周期时长，例如 2592000 = 30 天"
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
