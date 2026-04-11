import { Button } from "@mui/material";
import { WebAuthnP256 } from "ox";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "./api";

function App() {
  const { mutate, isPending } = useMutation({
    mutationFn: api.createAccount,
    onError: (err) => {
      console.error(err);
    },
  });

  const { data: supportedPermissions } = useQuery({
    queryKey: ["supportedPermissions"],
    queryFn: api.getSupportedPermissions,
  });

  async function handleClick() {
    const credential = await WebAuthnP256.createCredential({
      name: "wallet-user",
    });

    const publicKey = {
      prefix: credential.publicKey.prefix,
      x: `0x${credential.publicKey.x.toString(16).padStart(64, "0")}`,
      y: `0x${credential.publicKey.y.toString(16).padStart(64, "0")}`,
    };

    mutate({
      publicKeyX: publicKey.x,
      publicKeyY: publicKey.y,
      credentialId: credential.id,
    });
  }

  console.log("支持权限：", supportedPermissions);

  return (
    <>
      <Button variant="contained" onClick={handleClick}>
        Click Me
      </Button>
      {supportedPermissions && (
        <pre>{JSON.stringify(supportedPermissions, null, 2)}</pre>
      )}
    </>
  );
}

export default App;
