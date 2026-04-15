import {
  TextField,
  Stack,
  FormControlLabel,
  Checkbox,
  Typography,
} from "@mui/material";
import { useState, useEffect } from "react";

import type { ERC20TokenRevocationPermissionRequest } from "../../types/permissions";

type Props = {
  onChange: (req: ERC20TokenRevocationPermissionRequest) => void;
};

export function ERC20TokenRevocationForm({ onChange }: Props) {
  const now = Math.floor(Date.now() / 1000);

  const [expiry, setExpiry] = useState(String(now + 60 * 60 * 24 * 30));
  const [justification, setJustification] = useState(
    "This site needs to revoke your token approvals for safety.",
  );
  const [isAdjustmentAllowed, setIsAdjustmentAllowed] = useState(true);

  useEffect(() => {
    onChange({
      type: "erc20-token-revocation",
      startTime: null,
      expiry: expiry.trim() ? Number(expiry) : null,
      justification: justification || null,
      isAdjustmentAllowed,
    });
  }, [expiry, justification, isAdjustmentAllowed, onChange]);

  return (
    <Stack spacing={2}>
      <TextField
        label="Justification"
        multiline
        rows={3}
        fullWidth
        value={justification}
        onChange={(e) => setJustification(e.target.value)}
        helperText="说明撤销授权的原因"
      />
      <TextField
        label="Expiry (unix timestamp)"
        type="number"
        fullWidth
        value={expiry}
        onChange={(e) => setExpiry(e.target.value)}
        helperText="撤销权限的有效期"
      />
      <FormControlLabel
        control={
          <Checkbox
            checked={isAdjustmentAllowed}
            onChange={(e) => setIsAdjustmentAllowed(e.target.checked)}
          />
        }
        label={<Typography variant="body2">Allow Adjustments</Typography>}
      />
    </Stack>
  );
}
