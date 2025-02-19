import React, { useState, useEffect } from "react";
import { startAuthentication } from "@simplewebauthn/browser";
import {
  AuthenticationResponseJSON,
  PublicKeyCredentialCreationOptionsJSON,
} from "@simplewebauthn/types";
import { AuthResponse } from "../../utils/types";
import { Button } from "../input/Button";
import { Input } from "../input/Input";

export function PasswordLogin(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [notification, setNotification] = useState("");

  useEffect(() => {
    passkeyAutofill();
  });

  async function passkeyAutofill() {
    const response = await fetch(`/discoverable_login/begin`, {
      method: "POST",
      body: JSON.stringify({ email }),
      headers: {
        "Content-Type": "application/json",
      },
    });
    const credentialRequestOptions: {
      publicKey: PublicKeyCredentialCreationOptionsJSON;
    } = await response.json();
    let assertion: AuthenticationResponseJSON;
    try {
      assertion = await startAuthentication({
        optionsJSON: credentialRequestOptions.publicKey,
        useBrowserAutofill: true,
      });
    } catch (error) {
      if (error instanceof Error) {
        switch (error.name) {
          case "TypeError":
            setNotification("An account with that email does not exist.");
            break;
          case "AbortError":
            break;
          default:
            setNotification("An error occurred. Please try again.");
        }
      }
      return;
    }

    const verificationResponse = await fetch(`/discoverable_login/finish`, {
      method: "POST",
      body: JSON.stringify(assertion),
      headers: {
        "Content-Type": "application/json",
      },
    });

    const verificationJSON: AuthResponse = await verificationResponse.json();
    if (verificationJSON.status === "ok") {
      setNotification("Successfully logged in.");
    } else {
      setNotification("Login failed.");
    }
  }

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
