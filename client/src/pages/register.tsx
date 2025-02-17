import React, { useState } from "react";
import { RegistrationResponseJSON } from "@simplewebauthn/types";
import { startRegistration } from "@simplewebauthn/browser";

function Register(): React.ReactElement {
  const [username, setUsername] = useState("");
  const [notification, setNotification] = useState("");

  async function registerUser() {
    if (!username) {
      setNotification("Please enter your username");
      return;
    }

    try {
      const response = await fetch(`/register/begin/${username}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const credentialCreationOptions = await response.json();
      const registrationResponse: RegistrationResponseJSON =
        await startRegistration({
          optionsJSON: credentialCreationOptions.publicKey,
        });

      const verificationResponse = await fetch(`/register/finish/${username}`, {
        method: "POST",
        body: JSON.stringify(registrationResponse),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!verificationResponse.ok) {
        throw new Error(`HTTP error! status: ${verificationResponse.status}`);
      }
      const verificationResponseJSON = await verificationResponse.json();

      if (
        verificationResponseJSON &&
        verificationResponseJSON.status === "ok"
      ) {
        setNotification("Successfully registered.");
      } else {
        setNotification("Registration failed.");
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        switch (err.name) {
          case "InvalidStateError":
            setNotification("User already exists");
            break;
          default:
            setNotification(`Registration canceled or Failed due to an error`);
        }
      } else {
        setNotification("Registration failed due to an unknown error.");
      }
    }
  }

  return (
    <>
      <div className="header">
        <img
          src="/passkey_logo.png"
          width="60"
          height="60"
          alt="Passkey Logo"
        />
        <h1>Passkey Demo</h1>
      </div>
      <div id="notification">{notification}</div>
      <input
        type="text"
        name="username"
        id="username"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <button onClick={registerUser}>Register</button>
      <a className="link" href="/">
        Already have an account?
      </a>
    </>
  );
}

export default Register;
