import { PermissionType } from "@/enum";
import type { Hex } from "viem";

export const SUPPORTED_CHAINS = [
  { id: "0xaa36a7" as Hex, name: "Sepolia (Testnet)" },
  { id: "0x1" as Hex, name: "Ethereum Mainnet" },
  { id: "0x89" as Hex, name: "Polygon" },
  { id: "0xa4b1" as Hex, name: "Arbitrum One" },
  { id: "0xa" as Hex, name: "Optimism" },
  { id: "0x2105" as Hex, name: "Base" },
];

export const PERMISSION_TABS = [
  { value: PermissionType.NativeTokenStream, label: "Native Token Stream" },
  { value: PermissionType.ERC20TokenStream, label: "ERC20 Token Stream" },
  {
    value: PermissionType.NativeTokenPeriodic,
    label: "Native Token Periodic",
  },
  {
    value: PermissionType.ERC20TokenPeriodic,
    label: "ERC20 Token Periodic",
  },
  { value: PermissionType.ERC20TokenRevocation, label: "ERC20 Revocation" },
] as const;

// export type PermissionTabValue = (typeof PERMISSION_TABS)[number]["value"];
