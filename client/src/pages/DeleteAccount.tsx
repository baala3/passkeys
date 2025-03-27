import React, { useState } from "react";
import { Layout } from "../components/layout/Layout.tsx";
import { Button } from "../components/input/Button.tsx";
import { Heading } from "../components/layout/Heading.tsx";
import { loginPasskey } from "../hooks/webauth_api.tsx";
import { Notification } from "../components/layout/Notification";
import { Checkbox } from "../components/input/Checkbox";

export default function DeleteAccount(): React.ReactElement {
  const [confirm, setConfirm] = useState(false);
  const [notification, setNotification] = useState("");

  async function handleDeleteAccount() {
    if (!confirm) {
      setNotification("Confirm the below items");
      return;
    }

    loginPasskey(
      "",
      "delete_account",
      async () => {
        await fetch("/delete_account", {
          method: "DELETE",
        });
        window.location.reload();
      },
      (errorMessage) => setNotification(errorMessage)
    );
  }

  return (
    <Layout parent="/home">
      <Heading>Are you sure?</Heading>

      <Notification notification={notification} />

      <Checkbox
        checked={confirm}
        onChange={() => setConfirm(!confirm)}
        label="confirm that you have passkey to delete your account."
      />

      <br />
      <Button
        buttonText="Delete Account"
        onClickFunc={handleDeleteAccount}
        className="bg-red-500/10 hover:bg-red-500/20 text-red-500 border border-red-500/20 hover:border-red-500/30 hover:shadow-red-500/10"
      />
    </Layout>
  );
}
