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

export async function registerPasskey(
  email: string,
  context: string = "none",
  onSuccessCallback: () => void,
  onFailureCallback: (errorMessage: string) => void
) {
  if (context === "signup" && !isValidEmail(email)) {
    onFailureCallback("Please enter your email.");
    return;
  }

  const response = await fetch(`/register/begin?context=${context}`, {
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
      onFailureCallback(credentialCreationOptions.errorMessage);
      return;
    }

    registrationResponse = await startRegistration({
      optionsJSON: credentialCreationOptions.publicKey,
    });
  } catch (error) {
    if (error instanceof Error) {
      onFailureCallback("An error occurred. Please try again.");
    }
    return;
  }
  const verificationResponse = await fetch(
    `/register/finish?context=${context}`,
    {
      method: "POST",
      body: JSON.stringify(registrationResponse),
      headers: {
        "Content-Type": "application/json",
      },
    }
  );
  const verificationJSON: AuthResponse = await verificationResponse.json();

  if (verificationJSON.status === "ok") {
    onSuccessCallback();
  } else {
    onFailureCallback("Registration failed.");
  }
}

export async function loginPasskey(
  email: string,
  context: string = "none",
  onSuccessCallback: () => void,
  onFailureCallback: (errorMessage: string) => void
) {
  if (!isValidEmail(email)) {
    onFailureCallback("Please enter your email.");
    return;
  }

  const response = await fetch(`/login/begin?context=${context}`, {
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
          onFailureCallback(
            "There is no passkey associated with this account."
          );
          break;
        case "AbortError":
          break;
        default:
          onFailureCallback("An error occurred. Please try again.");
      }
    }
    return;
  }

  const verificationResponse = await fetch(`/login/finish?context=${context}`, {
    method: "POST",
    body: JSON.stringify(assertion),
    headers: {
      "Content-Type": "application/json",
    },
  });

  const verificationJSON: AuthResponse = await verificationResponse.json();
  if (verificationJSON.status === "ok") {
    onSuccessCallback();
  } else {
    onFailureCallback("Login failed.");
  }
}

export async function passkeyAutofill(
  email: string,
  context: string = "none",
  onSuccessCallback: () => void,
  onFailureCallback: (errorMessage: string) => void
) {
  const response = await fetch(`/discoverable_login/begin?context=${context}`, {
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
          onFailureCallback("An account with that email does not exist.");
          break;
        case "AbortError":
          break;
        default:
          onFailureCallback("An error occurred. Please try again.");
      }
    }
    return;
  }

  const verificationResponse = await fetch(
    `/discoverable_login/finish?context=${context}`,
    {
      method: "POST",
      body: JSON.stringify(assertion),
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  const verificationJSON: AuthResponse = await verificationResponse.json();
  if (verificationJSON.status === "ok") {
    onSuccessCallback();
  } else {
    onFailureCallback("Login failed.");
  }
}

export async function deletePasskey(
  credentialId: string,
  context: string = "none",
  onSuccessCallback: () => void,
  onFailureCallback: (errorMessage: string) => void
) {
  const response = await fetch(`/credentials?context=${context}`, {
    method: "DELETE",
    body: JSON.stringify({ credentialId: credentialId }),
    headers: {
      "Content-Type": "application/json",
    },
  });
  if (response.ok) {
    onSuccessCallback();
  } else {
    onFailureCallback("Failed to delete passkey.");
  }
}
