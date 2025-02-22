import React, { useState, useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { LinkButton } from "../components/input/LinkButton";
import { Button } from "../components/input/Button";
import { Heading } from "../components/layout/Heading";
import { HorizontalLine } from "../components/layout/HorizontalLine";
import { Passkey } from "../utils/types";
import { useNavigate } from "react-router-dom";
import { registerPasskey, deletePasskey } from "../hooks/webauth_api";
import { formatDate } from "../utils/shared";
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

  async function handleRegisterPasskey() {
    await registerPasskey("", "normal", navigate, setNotification);
  }

  async function handleDeletePasskey(credentialId: string) {
    await deletePasskey(credentialId);
  }

  return (
    <Layout>
      <Heading>Manage Passkeys</Heading>
      <div className="text-sm text-center min-h-8 font-normal text-blue-400">
        {notification}
      </div>
      <Button
        onClickFunc={handleRegisterPasskey}
        buttonText="Register a passkey"
      />
      <div className="font-light text-xs mt-2">
        A prompt will be displayed to confirm registration.
      </div>
      <HorizontalLine />
      {registeredPasskeys.map((passkey) => (
        <div>
          <div
            key={passkey.authenticator_metadata.name}
            className="grid grid-cols-2 gap-2 items-center"
          >
            <div>
              <div className="font-bold flex items-center gap-2">
                <img
                  src={passkey.authenticator_metadata.icon_light}
                  alt={passkey.authenticator_metadata.name}
                  className="w-6 h-6"
                />
                {passkey.authenticator_metadata.name}
              </div>
              <div className="font-light text-xs text-gray-400 mt-1">
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
