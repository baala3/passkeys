import React from "react";
import { Layout } from "../components/layout/Layout.tsx";
import { MenuItem } from "../components/navigation/MenuItem.tsx";
import { Button } from "../components/input/Button.tsx";

const MenuItems = [
  { title: "Manage Passkeys", link: "/passkeys" },
  { title: "Change email address", link: "/edit_email" },
  { title: "Change password", link: "/edit_password" },
  { title: "Lost passkey?", link: "#" },
  { title: "Delete Account", link: "/delete_account" },
];

export default function Homepage(): React.ReactElement {
  async function signOut() {
    await fetch("/logout", {
      method: "POST",
    });
    window.location.reload();
  }

  return (
    <Layout>
      <div className="space-y-3 mb-8">
        {MenuItems.map((item, index) => (
          <div
            key={index}
            className="group relative transition-all duration-200 hover:scale-[1.01]"
          >
            <MenuItem title={item.title} link={item.link} />
          </div>
        ))}
      </div>
      <div className="w-full">
        <Button
          onClickFunc={signOut}
          buttonText="Sign out"
          className="bg-red-500/10 hover:bg-red-500/20 text-red-500 border border-red-500/20 hover:border-red-500/30 hover:shadow-red-500/10"
        />
      </div>
    </Layout>
  );
}
