import React, { useState, useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { LinkButton } from "../components/input/LinkButton";
import { Button } from "../components/input/Button";
import { Heading } from "../components/layout/Heading";
import { HorizontalLine } from "../components/layout/HorizontalLine";
import { Passkey } from "../utils/types";

export default function ManagePasskeys(): React.ReactElement {
  const [registeredPasskeys, setRegisteredPasskeys] = useState<Passkey[]>([]);

  useEffect(() => {
    getPasskeys();
  }, []);

  async function getPasskeys() {
    const res = await fetch("/credentials");
    const passkeys = await res.json();
    setRegisteredPasskeys(passkeys);
  }

  return (
    <Layout>
      <Heading>Manage Passkeys</Heading>
      <Button
        onClickFunc={() => alert("Not implemented yet!")}
        buttonText="Register a passkey"
      />
      <div className="font-light text-xs mt-2">
        A prompt will be displayed to confirm registration.
      </div>
      <HorizontalLine />
      {registeredPasskeys.map((passkey) => (
        <div>
          <div key={passkey.aaguid} className="grid grid-cols-2 gap-4">
            <div>
              <div className="font-bold">{passkey.aaguid}</div>
              <div className="font-light text-xs text-gray-400">
                <p>Registered at: {passkey.created_at}</p>
                <p>Last used at: {passkey.updated_at}</p>
              </div>
            </div>
            <div>
              <LinkButton
                onClickFunc={() => alert("Not implemented yet!")}
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
