import type { Hex } from "viem";
import { PermissionType } from "@/enum";

// -------- 基础公共字段 --------
export type BasePermissionRequest = {
  justification: string | null;
  startTime: number | null;
  expiry: number | null;
  isAdjustmentAllowed: boolean;
};

// -------- 5 种权限类型 --------
export type NativeTokenStreamPermissionRequest = BasePermissionRequest & {
  type: typeof PermissionType.NativeTokenStream;
  initialAmount: bigint | null;
  amountPerSecond: bigint;
  maxAmount: bigint | null;
};

export type ERC20TokenStreamPermissionRequest = BasePermissionRequest & {
  type: typeof PermissionType.ERC20TokenStream;
  initialAmount: Hex | null;
  amountPerSecond: Hex;
  maxAmount: Hex | null;
  tokenAddress: Hex;
};

export type NativeTokenPeriodicPermissionRequest = BasePermissionRequest & {
  type: typeof PermissionType.NativeTokenPeriodic;
  periodAmount: Hex;
  periodDuration: number;
};

export type ERC20TokenPeriodicPermissionRequest = BasePermissionRequest & {
  type: typeof PermissionType.ERC20TokenPeriodic;
  periodAmount: Hex;
  periodDuration: number;
  tokenAddress: Hex;
};

export type ERC20TokenRevocationPermissionRequest = BasePermissionRequest & {
  type: typeof PermissionType.ERC20TokenRevocation;
};

export type PermissionRequest =
  | NativeTokenStreamPermissionRequest
  | ERC20TokenStreamPermissionRequest
  | NativeTokenPeriodicPermissionRequest
  | ERC20TokenPeriodicPermissionRequest
  | ERC20TokenRevocationPermissionRequest;

// -------- 发往后端的标准 ERC-7715 Payload --------
export type PermissionData = Record<string, unknown>;

export type ERC7715PermissionPayload = {
  chainId: Hex;
  to: Hex; // 委托账户地址
  expiry: number;
  isAdjustmentAllowed: boolean;
  permission: {
    type: string;
    data: string; // JSON.stringify(permissionData)
  };
};
