import APIClient from "./apiClient";

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

export default {
  createAccount,
  getSupportedPermissions,
};
