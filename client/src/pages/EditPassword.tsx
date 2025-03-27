import React, { useState } from "react";
import { Layout } from "../components/layout/Layout";
import { Heading } from "../components/layout/Heading";
import { Input } from "../components/input/Input";
import { Button } from "../components/input/Button";
import { loginPasskey } from "../hooks/webauth_api";
import { Notification } from "../components/layout/Notification";

export default function EditPassword(): React.ReactElement {
  const [newPassword, setNewPassword] = useState("");
  const [notification, setNotification] = useState("");

  async function handleChangePassword() {
    if (newPassword === "" || newPassword.length < 8) {
      setNotification("New password is short");
      return;
    }

    loginPasskey(
      "",
      "password_change",
      async () => {
        const response = await fetch("/change_password", {
          method: "POST",
          body: JSON.stringify({ password: newPassword }),
          headers: {
            "Content-Type": "application/json",
          },
        });
        if (response.ok) {
          setNotification("Password changed successfully");
        } else {
          setNotification("Failed to change password");
        }
      },
      (errorMessage) => setNotification(errorMessage)
    );
  }

  return (
    <Layout parent="/home">
      <Notification notification={notification} />
      <Heading>Edit Password</Heading>
      <p className="text-sm text-center font-normal text-gray-500 mb-4">
        Confirm that you have passkey to change your password.
      </p>
      <div className="space-y-6">
        <Input
          type="password"
          placeholder="New password"
          value={newPassword}
          onChange={setNewPassword}
        />

        <Button
          buttonText="Change password"
          onClickFunc={handleChangePassword}
        />
      </div>
    </Layout>
  );
}
