import React from "react";
import { Layout } from "../components/layout/Layout.tsx";
import { MenuItem } from "../components/navigation/MenuItem.tsx";
import { Button } from "../components/input/Button.tsx";

const MenuItems = [
  { title: "Change email address", link: "#" },
  { title: "Change password", link: "#" },
  { title: "Set up Two-Step Authentication", link: "#" },
  { title: "Manage Passkeys", link: "/passkeys" },
  // { title: "Delete Account", link: "#" },
];

export default function Homepage(): React.ReactElement {
  async function signOut() {
    await fetch("/logout", {
      method: "POST",
    });
    window.location.reload();
  }

  async function handleDeleteAccount() {
    await fetch("/delete_account", {
      method: "DELETE",
    });
    window.location.reload();
  }

  return (
    <Layout>
      <div className="mb-10">
        {MenuItems.map((item) => (
          <MenuItem title={item.title} link={item.link} />
        ))}
      </div>
      <Button onClickFunc={signOut} buttonText="Sign out" />
      <br />
      <Button onClickFunc={handleDeleteAccount} buttonText="Delete account" />
    </Layout>
  );
}
