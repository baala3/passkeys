import React, { useState } from "react";
import { Layout } from "../components/layout/Layout.tsx";
import { Button } from "../components/input/Button.tsx";
import { Heading } from "../components/layout/Heading.tsx";

export default function DeleteAccount(): React.ReactElement {
  const [confirm, setConfirm] = useState(false);
  const [notification, setNotification] = useState("");

  async function handleDeleteAccount() {
    if (!confirm) {
      setNotification("Confirm the below items");
      return;
    }

    await fetch("/delete_account", {
      method: "DELETE",
    });
    window.location.reload();
  }

  return (
    <Layout>
      <Heading>Are you sure?</Heading>

      <div className="text-sm text-center font-normal text-blue-400 mb-4">
        {notification}
      </div>

      <div className="text-center text-base text-gray-500 flex items-center justify-center leading-6">
        <input
          type="checkbox"
          className="mr-2 accent-[#027D9C]"
          checked={confirm}
          onChange={() => setConfirm(!confirm)}
        />
        <span>confirm that you have passkey to delete your account.</span>
      </div>

      <br />
      <Button buttonText="Delete Account" onClickFunc={handleDeleteAccount} />
    </Layout>
  );
}
