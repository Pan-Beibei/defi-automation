import APIClient from "./apiClient";
import type { ERC7715PermissionPayload } from "../types/permissions";

const apiClient = new APIClient();

interface CreateAccountRequest {
  publicKeyX: string;
  publicKeyY: string;
  credentialId: string;
}

const createAccount = (request: CreateAccountRequest) => {
  return apiClient.post({
    url: "/account/create",
    data: request,
  });
};

const getSupportedPermissions = () => {
  return apiClient.get({
    url: "/supported-permissions",
  });
};

const submitPermission = (payload: ERC7715PermissionPayload) => {
  return apiClient.post({
    url: "/permissions",
    data: payload,
  });
};

export default {
  createAccount,
  getSupportedPermissions,
  submitPermission,
};
