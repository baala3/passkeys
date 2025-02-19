import React, { useEffect } from "react";
import { Divider } from "../components/layout/Divider";
import { Link } from "../components/input/Link";
import { Layout } from "../components/layout/Layout";
import { PasswordSignUp } from "../components/form/PasswordSignUp";
import { PasskeySignUp } from "../components/form/PasskeySignUp";
import { useNavigate } from "react-router-dom";
import { isAuthenticated } from "../utils/shared";

export default function Register(): React.ReactElement {
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated()) {
      navigate("/home");
    }
  });

  return (
    <Layout>
      <PasswordSignUp />
      <Divider />
      <PasskeySignUp />

      <div className="mt-10">
        <Link href="/" linkText="Already have an account?" />
      </div>
    </Layout>
  );
}
