import {
  RegistrationResponseJSON,
  startRegistration,
  startAuthentication,
} from "@simplewebauthn/browser";
import {
  AuthenticationResponseJSON,
  PublicKeyCredentialCreationOptionsJSON,
} from "@simplewebauthn/types";
import { isValidEmail } from "../utils/shared";
import { AuthResponse } from "../utils/types";
import { NavigateFunction } from "react-router-dom";

export async function registerUser(
  email: string,
  navigate: NavigateFunction,
  setNotification: (notification: string) => void
) {
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

export async function loginUser(
  email: string,
  navigate: NavigateFunction,
  setNotification: (notification: string) => void
) {
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
  const credentialRequestOptions: {
    publicKey: PublicKeyCredentialCreationOptionsJSON;
  } = await response.json();
  let assertion: AuthenticationResponseJSON;
  try {
    assertion = await startAuthentication({
      optionsJSON: credentialRequestOptions.publicKey,
    });
  } catch (error) {
    if (error instanceof Error) {
      switch (error.name) {
        case "TypeError":
          setNotification("There is no passkey associated with this account.");
          break;
        case "AbortError":
          break;
        default:
          setNotification("An error occurred. Please try again.");
      }
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

  const verificationJSON: AuthResponse = await verificationResponse.json();
  if (verificationJSON.status === "ok") {
    setNotification("Successfully logged in.");
    navigate("/home");
  } else {
    setNotification("Login failed.");
  }
}

export async function passkeyAutofill(
  email: string,
  navigate: NavigateFunction,
  setNotification: (notification: string) => void
) {
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
    navigate("/home");
  } else {
    setNotification("Login failed.");
  }
}
