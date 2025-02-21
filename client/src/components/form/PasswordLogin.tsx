import React, { useState, useEffect } from "react";
import { AuthResponse } from "../../utils/types";
import { Button } from "../input/Button";
import { Input } from "../input/Input";
import { useNavigate } from "react-router-dom";
import { passkeyAutofill } from "../../hooks/webauth_api";

export function PasswordLogin(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();
  useEffect(() => {
    passkeyAutofill(email, navigate, setNotification);
  }, []);

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
      <h2 className="text-center text-xl font-bold leading-9 tracking-tight text-gray-900">
        Sign in with your password
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

        <Button onClickFunc={loginUser} buttonText="Sign in" />
      </div>
    </>
  );
}
