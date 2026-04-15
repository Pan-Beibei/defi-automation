import { PermissionType } from "@/enum";

export const SUPPORTED_CHAINS = [
  { id: 11155111, name: "Sepolia (Testnet)" },
  { id: 1, name: "Ethereum Mainnet" },
  { id: 137, name: "Polygon" },
  { id: 42161, name: "Arbitrum One" },
  { id: 10, name: "Optimism" },
  { id: 8453, name: "Base" },
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
