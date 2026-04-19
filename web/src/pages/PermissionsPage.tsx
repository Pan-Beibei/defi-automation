import {
  Box,
  Button,
  CircularProgress,
  Container,
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Tab,
  Tabs,
  TextField,
  Typography,
} from "@mui/material";
import { useCallback, useState } from "react";
import type { Hex } from "viem";
import {
  ERC20TokenPeriodicForm,
  ERC20TokenRevocationForm,
  ERC20TokenStreamForm,
  NativeTokenPeriodicForm,
  NativeTokenStreamForm,
} from "@/components/permissions";
import type {
  ERC7715PermissionPayload,
  PermissionRequest,
} from "@/types/permissions";
import { buildERC7715Payload, bigintReplacer } from "@/utils/permissionBuilder";
import api from "@/api";
import { useMutation } from "@tanstack/react-query";
import { PERMISSION_TABS, SUPPORTED_CHAINS } from "@/config";
import { PermissionType, type PermissionTypeValue } from "@/enum";

function PermissionsPage() {
  // 基础设置
  const [chainId, setChainId] = useState<Hex>("0xaa36a7");
  const [delegateTo, setDelegateTo] = useState<string>("");

  // 权限表单
  const [activeTab, setActiveTab] = useState<PermissionTypeValue>(
    PermissionType.NativeTokenStream,
  );
  const [permissionRequest, setPermissionRequest] =
    useState<PermissionRequest | null>(null);

  const mutation = useMutation({
    mutationFn: (payload: ERC7715PermissionPayload) =>
      api.submitPermission(payload),
  });

  // 每次表单字段变动触发
  const onFormChange = useCallback((req: PermissionRequest) => {
    setPermissionRequest(req);
  }, []);

  // 切换权限类型时重置状态
  const handleTabChange = (
    _: React.SyntheticEvent,
    newVal: PermissionTypeValue,
  ) => {
    setActiveTab(newVal);
    setPermissionRequest(null);
    mutation.reset();
  };

  // 提交到后端
  const handleSubmit = async () => {
    if (!permissionRequest) {
      mutation.reset(); // 可选：清空上次结果
      return;
    }
    if (!delegateTo.trim()) {
      return;
    }

    const payload = buildERC7715Payload(
      permissionRequest,
      chainId,
      delegateTo as Hex,
    );

    console.log("提交的数据：", payload);

    mutation.mutate(payload);
  };

  // 实时生成 payload 预览（try-catch 防止表单中间状态 BigInt 异常）
  let payloadPreview: ReturnType<typeof buildERC7715Payload> | null = null;
  try {
    if (permissionRequest) {
      payloadPreview = buildERC7715Payload(
        permissionRequest,
        chainId,
        (delegateTo || "0x") as Hex,
      );
    }
  } catch {
    payloadPreview = null;
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      {/* ── 标题 ── */}
      <Typography variant="h4" fontWeight={700} gutterBottom>
        ERC-7715 Permission Request
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
        配置权限参数后，点击「提交到后端」将请求发送至服务端处理。
      </Typography>

      {/* ── 基础设置 ── */}
      <Paper variant="outlined" sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" sx={{ mb: 2 }}>
          基础设置
        </Typography>
        <Box
          sx={{
            display: "flex",
            gap: 2,
            flexWrap: "wrap",
            alignItems: "flex-start",
          }}
        >
          <FormControl sx={{ minWidth: 220 }}>
            <InputLabel id="chain-label">Chain</InputLabel>
            <Select
              labelId="chain-label"
              label="Chain"
              value={chainId}
              onChange={(e) => setChainId(e.target.value as Hex)}
            >
              {SUPPORTED_CHAINS.map((chain) => (
                <MenuItem key={chain.id} value={chain.id}>
                  {chain.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <TextField
            label="Delegate To（委托目标地址）"
            value={delegateTo}
            onChange={(e) => setDelegateTo(e.target.value)}
            placeholder="0x..."
            sx={{ flexGrow: 1, minWidth: 280 }}
            helperText="接收委托权限的账户地址（EIP-7710 delegatee）"
          />
        </Box>
      </Paper>

      {/* ── 权限类型 + 表单 ── */}
      <Paper variant="outlined" sx={{ mb: 3 }}>
        <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
          <Tabs
            value={activeTab}
            onChange={handleTabChange}
            variant="scrollable"
            scrollButtons="auto"
          >
            {PERMISSION_TABS.map((tab) => (
              <Tab key={tab.value} label={tab.label} value={tab.value} />
            ))}
          </Tabs>
        </Box>

        <Box sx={{ p: 3 }}>
          {activeTab === PermissionType.NativeTokenStream && (
            <NativeTokenStreamForm onChange={onFormChange} />
          )}
          {activeTab === PermissionType.ERC20TokenStream && (
            <ERC20TokenStreamForm onChange={onFormChange} />
          )}
          {activeTab === PermissionType.NativeTokenPeriodic && (
            <NativeTokenPeriodicForm onChange={onFormChange} />
          )}
          {activeTab === PermissionType.ERC20TokenPeriodic && (
            <ERC20TokenPeriodicForm onChange={onFormChange} />
          )}
          {activeTab === PermissionType.ERC20TokenRevocation && (
            <ERC20TokenRevocationForm onChange={onFormChange} />
          )}
        </Box>
      </Paper>

      {/* ── Payload 预览 ── */}
      {payloadPreview && (
        <Paper variant="outlined" sx={{ p: 3, mb: 3, bgcolor: "grey.50" }}>
          <Typography variant="subtitle2" color="text.secondary" gutterBottom>
            Payload 预览（将 POST 到后端的 JSON 数据）
          </Typography>
          <Box
            component="pre"
            sx={{
              m: 0,
              fontSize: "0.78rem",
              overflowX: "auto",
              whiteSpace: "pre-wrap",
              wordBreak: "break-all",
              color: "text.primary",
            }}
          >
            {JSON.stringify(payloadPreview, bigintReplacer, 2)}
          </Box>
        </Paper>
      )}

      {/* ── 提交按钮 ── */}
      <Button
        variant="contained"
        size="large"
        fullWidth
        onClick={handleSubmit}
        disabled={mutation.isPending}
        startIcon={
          mutation.isPending ? (
            <CircularProgress size={18} color="inherit" />
          ) : undefined
        }
        sx={{ mb: 3 }}
      >
        {mutation.isPending ? "提交中..." : "提交到后端"}
      </Button>
    </Container>
  );
}

export default PermissionsPage;
