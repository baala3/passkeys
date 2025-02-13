// check if the browser supports WebAuthn
document.addEventListener("DOMContentLoaded", function () {
  if (!window.PublicKeyCredential) {
    showNotification(
      "WebAuthn is not supported on this browser. Please use a modern browser to use this demo."
    );
  }
});

// decode the base64url encoded value to a Uint8Array
function bufferDecode(value) {
  // Convert base64url to base64 by replacing "-" with "+" and "_" with "/"
  value = value.replace(/-/g, "+").replace(/_/g, "/");

  return Uint8Array.from(atob(value), (c) => c.charCodeAt(0));
}

function bufferEncode(value) {
  return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}

// register a new user
function registerUser() {
  const username = document.getElementById("username").value;
  if (!username) {
    showNotification("Please enter your username");
    return;
  }

  fetch(`/register/begin/${username}`)
    .then((res) => res.json())
    .then((credentialCreationOptions) => {
      // decode challenge and user id as they are base64url encoded
      credentialCreationOptions.publicKey.challenge = bufferDecode(
        credentialCreationOptions.publicKey.challenge
      );
      credentialCreationOptions.publicKey.user.id = bufferDecode(
        credentialCreationOptions.publicKey.user.id
      );

      // create the credential
      return navigator.credentials.create({
        publicKey: credentialCreationOptions.publicKey,
      });
    })
    .then((credential) => {
      console.log("success fetching navigator.credentials.create()");
      let attestationObject = credential.response.attestationObject;
      let clientDataJSON = credential.response.clientDataJSON;
      let rawId = credential.rawId;

      fetch(`/register/finish/${username}`, {
        method: "POST",
        body: JSON.stringify({
          id: credential.id,
          rawID: bufferEncode(rawId),
          type: credential.type,
          response: {
            attestationObject: bufferEncode(attestationObject),
            clientDataJSON: bufferEncode(clientDataJSON),
          },
        }),
        headers: {
          "Content-Type": "application/json",
        },
      });
    })
    .then(() => {
      console.log("passkey registered successfully");
      showNotification("Successfully registered.");
    })
    .catch((err) => {
      console.error("Error:", err);
      showNotification("Registration failed.");
    });
}

// login a user
function loginUser() {
  const username = document.getElementById("username").value;
  if (!username) {
    showNotification("Please enter your username");
    return;
  }

  fetch(`/login/begin/${username}`)
    .then((res) => res.json())
    .then((credentialRequestOptions) => {
      credentialRequestOptions.publicKey.challenge = bufferDecode(
        credentialRequestOptions.publicKey.challenge
      );

      credentialRequestOptions.publicKey.allowCredentials.forEach(
        (credential) => {
          credential.id = bufferDecode(credential.id);
        }
      );

      return navigator.credentials.get({
        publicKey: credentialRequestOptions.publicKey,
      });
    })
    .then((assertion) => {
      console.log("success fetching navigator.credentials.get()");
      let authData = assertion.response.authenticatorData;
      let clientDataJSON = assertion.response.clientDataJSON;
      let rawId = assertion.rawId;
      let sig = assertion.response.signature;
      let userHandle = assertion.response.userHandle;

      fetch(`/login/finish/${username}`, {
        method: "POST",
        body: JSON.stringify({
          id: assertion.id,
          rawID: bufferEncode(rawId),
          type: assertion.type,
          response: {
            authenticatorData: bufferEncode(authData),
            clientDataJSON: bufferEncode(clientDataJSON),
            signature: bufferEncode(sig),
            userHandle: bufferEncode(userHandle),
          },
        }),
        headers: {
          "Content-Type": "application/json",
        },
      });
    })
    .then(() => {
      console.log("sign in with passkey successful");
      showNotification("Successfully logged in.");
    })
    .catch((err) => {
      console.error("Error:", err);
      showNotification("Login failed.");
    });
}

function showNotification(message) {
  document.getElementById("notification").innerHTML = message;
  setTimeout(() => {
    document.getElementById("notification").innerHTML = "";
  }, 5000);
}
