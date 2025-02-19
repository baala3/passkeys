import React, { useState } from "react";
import { Button } from "../input/Button";
import { Input } from "../input/Input";

export function PasswordSignUp(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [notification, setNotification] = useState("");

  return (
    <>
      <h2 className="text-center text-xl font-bold leading-9 tracking-tight text-gray-900">
        Sign up using a password
      </h2>
      <div className="space-y-6">
        <div className="text-sm text-center min-h-5 font-normal text-blue-400">
          {notification}
        </div>
        <Input
          type="email"
          placeholder="Email"
          value={email}
          onChange={setEmail}
        />
        <Input
          type="password"
          placeholder="Password"
          value={password}
          onChange={setPassword}
        />

        <Button
          onClickFunc={() => setNotification("Not yet implemented")}
          buttonText="Sign up"
        />
      </div>
    </>
  );
}
