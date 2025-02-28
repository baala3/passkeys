import React, { useState } from "react";
import { Layout } from "../components/layout/Layout";
import { Heading } from "../components/layout/Heading";
import { Input } from "../components/input/Input";
import { Button } from "../components/input/Button";
import { isValidEmail } from "../utils/shared";
import { loginPasskey } from "../hooks/webauth_api";
import { Notification } from "../components/layout/Notification";

export default function EditEmail(): React.ReactElement {
  const [newEmail, setNewEmail] = useState("");
  const [notification, setNotification] = useState("");

  async function handleChangeEmail() {
    if (newEmail === "" || !isValidEmail(newEmail)) {
      setNotification("New email is invalid");
      return;
    }

    loginPasskey(
      "",
      "email_change",
      async () => {
        const response = await fetch("/change_email", {
          method: "POST",
          body: JSON.stringify({ email: newEmail }),
          headers: {
            "Content-Type": "application/json",
          },
        });
        if (response.ok) {
          setNotification("Email changed successfully");
        } else {
          setNotification("Failed to change email");
        }
      },
      (errorMessage) => setNotification(errorMessage)
    );
  }

  return (
    <Layout>
      <Notification notification={notification} />
      <Heading>Edit Email</Heading>
      <p className="text-sm text-center font-normal text-gray-500 mb-4">
        Confirm that you have passkey to change your email.
      </p>
      <div className="space-y-6">
        <Input
          type="email"
          placeholder="New email"
          value={newEmail}
          onChange={setNewEmail}
        />

        <Button buttonText="Change email" onClickFunc={handleChangeEmail} />
      </div>
    </Layout>
  );
}
