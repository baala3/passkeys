import React from "react";
import { Layout } from "../components/layout/Layout";
import { MenuItem } from "../components/navigation/MenuItem.tsx";
import { Button } from "../components/input/Button.tsx";

const MenuItems = [
  { title: "Change email address", link: "#" },
  { title: "Change password", link: "#" },
  { title: "Set up Two-Step Authentication", link: "#" },
  { title: "Manage Passkeys", link: "#" },
  { title: "Delete Account", link: "#" },
];

export default function Homepage(): React.ReactElement {
  return (
    <Layout>
      <div className="mb-10">
        {MenuItems.map((item) => (
          <MenuItem title={item.title} link={item.link} />
        ))}
      </div>
      <Button
        onClickFunc={() => alert("Not yet implemented!")}
        buttonText="Sign out"
      />
    </Layout>
  );
}
