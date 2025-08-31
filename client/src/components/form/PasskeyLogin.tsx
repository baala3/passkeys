import React, { useState, useEffect } from "react";
import { Button } from "../input/Button";
import { Input } from "../input/Input";
import { useNavigate } from "react-router-dom";
import { loginPasskey } from "../../hooks/webauth_api";
import { Notification } from "../layout/Notification";
import { SubHeading } from "../layout/SubHeading";
import { Checkbox } from "../input/Checkbox";
import { passkeyAutofill } from "../../hooks/webauth_api";

export function PasskeyLogin(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();
  const [isAutofill, setIsAutofill] = useState(localStorage.getItem("isAutofill") === "true");

  async function handleLoginPasskey() {
    await loginPasskey(
      email,
      "signin",
      async () => navigate("/home"),
      (errorMessage) => setNotification(errorMessage)
    );
  }

  useEffect(() => {
    isAutofill && passkeyAutofill(
      email,
      "login",
      () => { console.log("success from autofill"); navigate("/home") },
      (errorMessage) => { console.log("error from autofill", errorMessage); setNotification(errorMessage) }
    );
  }, []);

  function handleAutofillChange() {
    localStorage.setItem("isAutofill", (!isAutofill).toString());
    setIsAutofill(!isAutofill);
    window.location.reload();
  }

  return (
    <>
      <SubHeading>Sign in with passkey</SubHeading>

      <div className="space-y-6">
        <Notification notification={notification} />

        <Checkbox
          checked={isAutofill}
          onChange={handleAutofillChange}
          label="use passkey autofill"
        />

        <Input
          type="email"
          placeholder={isAutofill ? "Passkey Autofill" : "Normal Passkey Flow"}
          value={email}
          onChange={setEmail}
          autoComplete={isAutofill ? "webauthn" : "off"}
        />

        <Button onClickFunc={handleLoginPasskey} buttonText="Sign in" />
      </div>
    </>
  );
}
