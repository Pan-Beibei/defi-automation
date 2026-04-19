import { toHex } from "viem";
import type { Hex } from "viem";

import type {
  PermissionRequest,
  ERC7715PermissionPayload,
  PermissionData,
} from "../types/permissions";

/**
 * 将前端表单数据打包成 ERC-7715 标准 payload，发往后端
 */
export function buildERC7715Payload(
  request: PermissionRequest,
  chainId: Hex,
  delegateTo: Hex,
): ERC7715PermissionPayload {
  const { type, expiry, isAdjustmentAllowed, ...rest } = request;

  let permissionData: PermissionData;

  switch (type) {
    case "native-token-stream": {
      const { initialAmount, amountPerSecond, maxAmount, startTime } = rest as {
        initialAmount: bigint | null;
        amountPerSecond: bigint;
        maxAmount: bigint | null;
        startTime: number | null;
      };
      permissionData = {
        initialAmount: initialAmount != null ? toHex(initialAmount) : null,
        amountPerSecond: toHex(amountPerSecond),
        maxAmount: maxAmount != null ? toHex(maxAmount) : null,
        startTime,
      };
      break;
    }

    case "erc20-token-stream": {
      const {
        initialAmount,
        amountPerSecond,
        maxAmount,
        tokenAddress,
        startTime,
      } = rest as {
        initialAmount: Hex | null;
        amountPerSecond: Hex;
        maxAmount: Hex | null;
        tokenAddress: Hex;
        startTime: number | null;
      };
      permissionData = {
        tokenAddress,
        initialAmount,
        amountPerSecond,
        maxAmount,
        startTime,
      };
      break;
    }

    case "native-token-periodic": {
      const { periodAmount, periodDuration, startTime } = rest as {
        periodAmount: Hex;
        periodDuration: number;
        startTime: number | null;
      };
      permissionData = { periodAmount, periodDuration, startTime };
      break;
    }

    case "erc20-token-periodic": {
      const { periodAmount, periodDuration, tokenAddress, startTime } =
        rest as {
          periodAmount: Hex;
          periodDuration: number;
          tokenAddress: Hex;
          startTime: number | null;
        };
      permissionData = {
        tokenAddress,
        periodAmount,
        periodDuration,
        startTime,
      };
      break;
    }

    case "erc20-token-revocation": {
      permissionData = {};
      break;
    }

    default:
      throw new Error(`Unknown permission type: ${String(type)}`);
  }

  return {
    chainId,
    to: delegateTo,
    expiry: expiry ?? 0,
    isAdjustmentAllowed,
    permission: { type, data: JSON.stringify(permissionData) },
  };
}

// bigint 无法直接 JSON.stringify，需要自定义 replacer
export function bigintReplacer(_key: string, value: unknown): unknown {
  return typeof value === "bigint" ? value.toString() : value;
}
