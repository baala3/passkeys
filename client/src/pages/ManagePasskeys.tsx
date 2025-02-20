import React, { useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { useNavigate } from "react-router-dom";
import { isAuthenticated } from "../utils/shared.ts";
import { LinkButton } from "../components/input/LinkButton";
import { Button } from "../components/input/Button";
import { Heading } from "../components/layout/Heading";
import { HorizontalLine } from "../components/layout/HorizontalLine";

export default function ManagePasskeys(): React.ReactElement {
  const navigate = useNavigate();

  useEffect(() => {
    if (!isAuthenticated()) {
      navigate("/");
    }
  }, []);

  // TODO: Replace with actual passkeys
  const registeredPasskeys = [
    {
      name: "Chrome on Mac",
      registeredAt: "July 19, 2024",
      lastUsedAt: "July 19, 2024",
    },
    {
      name: "Edge on Windows",
      registeredAt: "July 17, 2024",
      lastUsedAt: "July 19, 2024",
    },
  ];

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
          <div key={passkey.name} className="grid grid-cols-2 gap-4">
            <div>
              <div className="font-bold">{passkey.name}</div>
              <div className="font-light text-xs text-gray-400">
                <p>Registered at: {passkey.registeredAt}</p>
                <p>Last used at: {passkey.lastUsedAt}</p>
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
