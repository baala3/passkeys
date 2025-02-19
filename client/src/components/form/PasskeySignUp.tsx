import React, { useState } from "react";
import { startRegistration } from "@simplewebauthn/browser";
import {
  RegistrationResponseJSON,
  PublicKeyCredentialCreationOptionsJSON,
} from "@simplewebauthn/types";
import { Button } from "../input/Button";
import { Input } from "../input/Input";
import { isValidEmail } from "../../utils/shared";
import { AuthResponse } from "../../utils/types";
import { useNavigate } from "react-router-dom";

export function PasskeySignUp(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();

  async function registerUser() {
    if (!isValidEmail(email)) {
      setNotification("Please enter your email.");
      return;
    }

    const response = await fetch(`/register/begin`, {
      method: "POST",
      body: JSON.stringify({ email }),
      headers: {
        "Content-Type": "application/json",
      },
    });
    let registrationResponse: RegistrationResponseJSON;
    try {
      const credentialCreationOptions: {
        publicKey: PublicKeyCredentialCreationOptionsJSON;
      } & AuthResponse = await response.json();

      if (credentialCreationOptions.status === "error") {
        setNotification(credentialCreationOptions.errorMessage);
        return;
      }

      registrationResponse = await startRegistration({
        optionsJSON: credentialCreationOptions.publicKey,
      });
    } catch (error) {
      if (error instanceof Error) {
        setNotification("An error occurred. Please try again.");
      }
      return;
    }
    const verificationResponse = await fetch(`/register/finish`, {
      method: "POST",
      body: JSON.stringify(registrationResponse),
      headers: {
        "Content-Type": "application/json",
      },
    });
    const verificationJSON: AuthResponse = await verificationResponse.json();

    if (verificationJSON.status === "ok") {
      setNotification("Successfully registered.");
      navigate("/home");
    } else {
      setNotification("Registration failed.");
    }
  }

  return (
    <>
      <h2 className="text-center text-xl font-bold leading-9 tracking-tight text-gray-900">
        Create a new account with passkey
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

        <Button onClickFunc={registerUser} buttonText="Sign up" />
      </div>
    </>
  );
}
