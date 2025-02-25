import React, { useState } from "react";
import { Layout } from "../components/layout/Layout";
import { Heading } from "../components/layout/Heading";
import { Input } from "../components/input/Input";
import { Button } from "../components/input/Button";
import { isValidEmail } from "../utils/shared";
import { loginPasskey } from "../hooks/webauth_api";
export default function EditEmail(): React.ReactElement {
  const [currentEmail, setCurrentEmail] = useState("");
  const [newEmail, setNewEmail] = useState("");
  const [notification, setNotification] = useState("");

  async function handleChangeEmail() {
    if (!emailValidations()) {
      return;
    }

    loginPasskey(
      "",
      "email_change",
      async () => {
        const response = await fetch("/change_email", {
          method: "POST",
          body: JSON.stringify({ currentEmail, newEmail }),
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

  function emailValidations(): boolean {
    if (currentEmail === "") {
      setNotification("Current email is required");
      return false;
    }
    if (newEmail === "") {
      setNotification("New email is required");
      return false;
    }
    if (!isValidEmail(currentEmail) || !isValidEmail(newEmail)) {
      setNotification("Invalid email entered");
      return false;
    }
    if (currentEmail === newEmail) {
      setNotification("Current and new email cannot be the same");
      return false;
    }
    return true;
  }

  return (
    <Layout>
      <div className="text-sm text-center font-normal text-blue-400 mb-4">
        {notification}
      </div>
      <Heading>Edit Email</Heading>
      <p className="text-sm text-center font-normal text-gray-500 mb-4">
        Confirm that you have passkey to change your email.
      </p>
      <div className="space-y-6">
        <Input
          type="email"
          placeholder="Current email"
          value={currentEmail}
          onChange={setCurrentEmail}
        />
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
