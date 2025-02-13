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

function showNotification(message) {
  document.getElementById("notification").innerHTML = message;
  setTimeout(() => {
    document.getElementById("notification").innerHTML = "";
  }, 5000);
}
