export const PermissionType = {
  NativeTokenStream: "native-token-stream",
  ERC20TokenStream: "erc20-token-stream",
  NativeTokenPeriodic: "native-token-periodic",
  ERC20TokenPeriodic: "erc20-token-periodic",
  ERC20TokenRevocation: "erc20-token-revocation",
} as const;

export type PermissionTypeValue =
  (typeof PermissionType)[keyof typeof PermissionType];
