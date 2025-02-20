import React, { useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { useNavigate } from "react-router-dom";
import { isAuthenticated } from "../utils/shared.ts";

export default function ManagePasskeys(): React.ReactElement {
  const navigate = useNavigate();

  useEffect(() => {
    if (!isAuthenticated()) {
      navigate("/");
    }
  }, []);

  return <Layout>Manage Passkey Page</Layout>;
}
