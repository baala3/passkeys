// check if the browser supports WebAuthn
document.addEventListener("DOMContentLoaded", function () {
  if (!window.PublicKeyCredential) {
    showNotification(
      "WebAuthn is not supported on this browser. Please use a modern browser to use this demo."
    );
  }
});

function showNotification(message) {
  document.getElementById("notification").innerHTML = message;
  setTimeout(() => {
    document.getElementById("notification").innerHTML = "";
  }, 5000);
}
