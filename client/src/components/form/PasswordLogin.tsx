import React, { useState } from "react";
import { AuthResponse } from "../../utils/types";
import { Button } from "../input/Button";
import { Input } from "../input/Input";
import { useNavigate } from "react-router-dom";
import { Notification } from "../layout/Notification";
import { SubHeading } from "../layout/SubHeading";

export function PasswordLogin(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();

  async function loginUser() {
    if (email === "") {
      setNotification("Please enter your email.");
      return;
    }

    if (password === "") {
      setNotification("Please enter your password.");
      return;
    }

    const response = await fetch(`/login/password`, {
      method: "POST",
      body: JSON.stringify({ email, password }),
      headers: {
        "Content-Type": "application/json",
      },
    });
    const loginJSON: AuthResponse = await response.json();
    if (loginJSON.status === "ok") {
      setNotification("Successfully logged in.");
      navigate("/home");
    } else {
      setNotification(loginJSON.errorMessage);
    }
  }

  return (
    <>
      <SubHeading>Sign in with your password</SubHeading>
      <div className="space-y-6">
        <Notification notification={notification} />
        <Input
          type="email"
          placeholder="email"
          value={email}
          onChange={setEmail}
        />
        <Input
          type="password"
          placeholder="password"
          value={password}
          onChange={setPassword}
        />

        <Button onClickFunc={loginUser} buttonText="Sign in" />
      </div>
    </>
  );
}
