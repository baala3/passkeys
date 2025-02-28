import React, { useState } from "react";
import { Button } from "../input/Button";
import { Input } from "../input/Input";
import { isValidEmail } from "../../utils/shared";
import { AuthResponse } from "../../utils/types";
import { useNavigate } from "react-router-dom";
import { Notification } from "../layout/Notification";
import { SubHeading } from "../layout/SubHeading";

export function PasswordSignUp(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();

  async function registerUser() {
    if (!isValidEmail(email)) {
      setNotification("Please enter a valid email.");
      return;
    }

    if (password.length < 8) {
      setNotification("Password must be at least 8 characters long.");
      return;
    }

    const response = await fetch(`/register/password`, {
      method: "POST",
      body: JSON.stringify({ email, password }),
      headers: {
        "Content-Type": "application/json",
      },
    });

    const registrationJSON: AuthResponse = await response.json();
    if (registrationJSON.status === "ok") {
      setNotification("Successfully registered.");
      navigate("/home");
    } else {
      setNotification(registrationJSON.errorMessage);
    }
  }

  return (
    <>
      <SubHeading>Sign up using a password</SubHeading>
      <div className="space-y-6">
        <Notification notification={notification} />
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

        <Button onClickFunc={registerUser} buttonText="Sign up" />
      </div>
    </>
  );
}
