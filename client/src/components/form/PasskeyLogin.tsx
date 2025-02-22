import React, { useState } from "react";
import { Button } from "../input/Button";
import { Input } from "../input/Input";
import { useNavigate } from "react-router-dom";
import { loginPasskey } from "../../hooks/webauth_api";

export function PasskeyLogin(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();

  async function handleLoginPasskey() {
    await loginPasskey(email, navigate, setNotification);
  }

  return (
    <>
      <h2 className="text-center text-xl font-bold leading-9 tracking-tight text-gray-900">
        Sign in with passkey
      </h2>

      <div className="space-y-6">
        <div className="text-sm text-center min-h-8 font-normal text-blue-400">
          {notification}
        </div>

        <Input
          type="email"
          placeholder="Email"
          value={email}
          onChange={setEmail}
        />

        <Button onClickFunc={handleLoginPasskey} buttonText="Sign in" />
      </div>
    </>
  );
}
