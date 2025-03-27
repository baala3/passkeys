import React, { useState } from "react";
import { Button } from "../input/Button";
import { Input } from "../input/Input";
import { useNavigate } from "react-router-dom";
import { loginPasskey } from "../../hooks/webauth_api";
import { Notification } from "../layout/Notification";
import { SubHeading } from "../layout/SubHeading";

export function PasskeyLogin(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();

  async function handleLoginPasskey() {
    await loginPasskey(
      email,
      "signin",
      async () => navigate("/home"),
      (errorMessage) => setNotification(errorMessage)
    );
  }

  return (
    <>
      <SubHeading>Sign in with passkey</SubHeading>

      <div className="space-y-6">
        <Notification notification={notification} />

        <Input
          type="email"
          placeholder="email"
          value={email}
          onChange={setEmail}
        />

        <Button onClickFunc={handleLoginPasskey} buttonText="Sign in" />
      </div>
    </>
  );
}
