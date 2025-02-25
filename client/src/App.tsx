import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Login from "./pages/Login.tsx";
import Register from "./pages/Register.tsx";
import Homepage from "./pages/Homepage.tsx";
import ManagePasskeys from "./pages/ManagePasskeys.tsx";
import DeleteAccount from "./pages/DeleteAccount.tsx";
import EditEmail from "./pages/EditEmail.tsx";
import EditPassword from "./pages/EditPassword.tsx";

function App(): React.ReactElement {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/sign-up" element={<Register />} />
        <Route path="/home" element={<Homepage />} />
        <Route path="/passkeys" element={<ManagePasskeys />} />
        <Route path="/delete_account" element={<DeleteAccount />} />
        <Route path="/edit_email" element={<EditEmail />} />
        <Route path="/edit_password" element={<EditPassword />} />
      </Routes>
    </Router>
  );
}

export default App;
