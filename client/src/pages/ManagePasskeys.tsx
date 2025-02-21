import React, { useState, useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { LinkButton } from "../components/input/LinkButton";
import { Button } from "../components/input/Button";
import { Heading } from "../components/layout/Heading";
import { HorizontalLine } from "../components/layout/HorizontalLine";
import { Passkey } from "../utils/types";
import DateObject from "react-date-object";
import { useNavigate } from "react-router-dom";
import { registerUser } from "../hooks/webauth_api";

export default function ManagePasskeys(): React.ReactElement {
  const [registeredPasskeys, setRegisteredPasskeys] = useState<Passkey[]>([]);
  const [notification, setNotification] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    getPasskeys();
  }, []);

  async function getPasskeys() {
    const res = await fetch("/credentials");
    const passkeys = (await res.json()) || [];
    setRegisteredPasskeys(passkeys);
  }

  function formatDate(date: string) {
    const dateObject = new DateObject({ date: new Date(date).toISOString() });
    return dateObject.format("DD MMM YY, hh:mm a");
  }

  async function handleRegisterUser() {
    await registerUser("", "none", navigate, setNotification);
  }

  async function handleDeletePasskey(credentialId: string) {
    const response = await fetch(`/credentials`, {
      method: "DELETE",
      body: JSON.stringify({ credentialId: credentialId }),
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (response.ok) {
      window.location.reload();
    }
  }

  return (
    <Layout>
      <Heading>Manage Passkeys</Heading>
      <div className="text-sm text-center min-h-5 font-normal text-blue-400">
        {notification}
      </div>
      <Button
        onClickFunc={() => handleRegisterUser()}
        buttonText="Register a passkey"
      />
      <div className="font-light text-xs mt-2">
        A prompt will be displayed to confirm registration.
      </div>
      <HorizontalLine />
      {registeredPasskeys.map((passkey) => (
        <div>
          <div key={passkey.credential_id} className="grid grid-cols-2 gap-4">
            <div>
              <div className="font-bold">{passkey.credential_id}</div>
              <div className="font-light text-xs text-gray-400">
                <p>Registered: {formatDate(passkey.created_at)}</p>
                <p>Last-Used: {formatDate(passkey.updated_at)}</p>
              </div>
            </div>
            <div>
              <LinkButton
                onClickFunc={() => handleDeletePasskey(passkey.credential_id)}
                buttonText="Delete"
              />
            </div>
          </div>
          <HorizontalLine />
        </div>
      ))}
    </Layout>
  );
}
