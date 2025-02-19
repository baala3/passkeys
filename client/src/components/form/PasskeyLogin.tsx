import React, { useState } from "react";
import { startAuthentication } from "@simplewebauthn/browser";
import { AuthenticationResponseJSON } from "@simplewebauthn/types";
import { Button } from "../input/Button";
import { Input } from "../input/Input";

export function PasskeyLogin(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [notification, setNotification] = useState("");

  function isValidEmail(email: string): boolean {
    const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/g;
    return email !== "" && emailRegex.test(email);
  }

  async function loginUser() {
    if (!isValidEmail(email)) {
      setNotification("Please enter your email.");
      return;
    }

    const response = await fetch(`/login/begin`, {
      method: "POST",
      body: JSON.stringify({ email }),
      headers: {
        "Content-Type": "application/json",
      },
    });
    const credentialRequestOptions = await response.json();
    let assertion: AuthenticationResponseJSON;
    try {
      assertion = await startAuthentication({
        optionsJSON: credentialRequestOptions.publicKey,
      });
    } catch (error: any) {
      switch (error.name) {
        case "TypeError":
          setNotification("An account with that email does not exist.");
          break;
        case "AbortError":
          break;
        default:
          setNotification("An error occurred. Please try again.");
      }
      return;
    }

    const verificationResponse = await fetch(`/login/finish`, {
      method: "POST",
      body: JSON.stringify(assertion),
      headers: {
        "Content-Type": "application/json",
      },
    });

    const verificationJSON = await verificationResponse.json();
    if (verificationJSON && verificationJSON.status === "ok") {
      setNotification("Successfully logged in.");
    } else {
      setNotification("Login failed.");
    }
  }

  return (
    <>
      <h2 className="text-center text-xl font-bold leading-9 tracking-tight text-gray-900">
        Sign in with passkey
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

        <Button onClickFunc={loginUser} buttonText="Sign in" />
      </div>
    </>
  );
}
