import React, { useState } from "react";
import { Button } from "../input/Button";
import { Input } from "../input/Input";
import { registerPasskey } from "../../hooks/webauth_api";
import { useNavigate } from "react-router-dom";
import { Notification } from "../layout/Notification";
import { SubHeading } from "../layout/SubHeading";

export function PasskeySignUp(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();

  async function handleRegisterPasskey() {
    await registerPasskey(
      email,
      "signup",
      () => navigate("/home"),
      (errorMessage) => setNotification(errorMessage)
    );
  }

  return (
    <>
      <SubHeading>Create a new account with passkey</SubHeading>

      <div className="space-y-6">
        <Notification notification={notification} />

        <Input
          type="email"
          placeholder="Email"
          value={email}
          onChange={setEmail}
        />

        <Button onClickFunc={handleRegisterPasskey} buttonText="Sign up" />
      </div>
    </>
  );
}
