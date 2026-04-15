import {
  TextField,
  FormControlLabel,
  Checkbox,
  Stack,
  Typography,
} from "@mui/material";

type BasePermissionFieldsProps = {
  justification: string;
  startTime: string; // unix timestamp string，空字符串表示 null
  expiry: string; // unix timestamp string
  isAdjustmentAllowed: boolean;
  onJustificationChange: (v: string) => void;
  onStartTimeChange: (v: string) => void;
  onExpiryChange: (v: string) => void;
  onIsAdjustmentAllowedChange: (v: boolean) => void;
};

export function BasePermissionFields({
  justification,
  startTime,
  expiry,
  isAdjustmentAllowed,
  onJustificationChange,
  onStartTimeChange,
  onExpiryChange,
  onIsAdjustmentAllowedChange,
}: BasePermissionFieldsProps) {
  return (
    <>
      <TextField
        label="Justification"
        multiline
        rows={3}
        fullWidth
        value={justification}
        onChange={(e) => onJustificationChange(e.target.value)}
        helperText="说明请求此权限的原因"
      />
      <Stack direction="row" spacing={2}>
        <TextField
          label="Start Time (unix timestamp)"
          type="number"
          fullWidth
          value={startTime}
          onChange={(e) => onStartTimeChange(e.target.value)}
          helperText="留空表示立即生效"
        />
        <TextField
          label="Expiry (unix timestamp)"
          type="number"
          fullWidth
          value={expiry}
          onChange={(e) => onExpiryChange(e.target.value)}
          helperText="权限过期时间"
        />
      </Stack>
      <FormControlLabel
        control={
          <Checkbox
            checked={isAdjustmentAllowed}
            onChange={(e) => onIsAdjustmentAllowedChange(e.target.checked)}
          />
        }
        label={
          <Typography variant="body2">
            Allow Adjustments（允许调整权限参数）
          </Typography>
        }
      />
    </>
  );
}
