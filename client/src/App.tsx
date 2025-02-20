import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Login from "./pages/login.tsx";
import Register from "./pages/register.tsx";
import Homepage from "./pages/homepage.tsx";
import ManagePasskeys from "./pages/ManagePasskeys.tsx";

function App(): React.ReactElement {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/sign-up" element={<Register />} />
        <Route path="/home" element={<Homepage />} />
        <Route path="/passkeys" element={<ManagePasskeys />} />
      </Routes>
    </Router>
  );
}

export default App;
