import React, { useState } from "react";
import { startAuthentication } from "@simplewebauthn/browser";
import { AuthenticationResponseJSON } from "@simplewebauthn/types";

function Login(): React.ReactElement {
  const [username, setUsername] = useState("");
  const [notification, setNotification] = useState("");

  async function loginUser() {
    if (!username) {
      setNotification("Please enter your username");
      return;
    }

    try {
      const response = await fetch(`/login/begin/${username}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const credentialRequestOptions = await response.json();

      const assertion: AuthenticationResponseJSON = await startAuthentication({
        optionsJSON: credentialRequestOptions.publicKey,
      });

      const verificationResponse = await fetch(`/login/finish/${username}`, {
        method: "POST",
        body: JSON.stringify(assertion),
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
        setNotification("Successfully logged in.");
      } else {
        setNotification("Login failed.");
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        switch (err.name) {
          case "TypeError":
            setNotification("User does not exist");
            break;
          default:
            setNotification("Login canceled or Failed due to an error.");
        }
      } else {
        setNotification("Login failed due to an unknown error.");
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
        id="username"
        name="username"
        placeholder="Enter username"
        autoComplete="username webauthn"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <button onClick={loginUser}>Login</button>
      <a className="link" href="/sign-up">
        Don't have an account?
      </a>
    </>
  );
}

export default Login;
