import React, { useState } from "react";
import { startAuthentication } from "@simplewebauthn/browser";
import { AuthenticationResponseJSON } from "@simplewebauthn/types";

function Login(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [notification, setNotification] = useState("");

  async function loginUser() {
    const isValidEmail = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/g;

    if (!email || !isValidEmail.test(email)) {
      setNotification("Please enter your email");
      return;
    }

    try {
      const response = await fetch(`/login/begin`, {
        method: "POST",
        body: JSON.stringify({ email }),
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const credentialRequestOptions = await response.json();

      const assertion: AuthenticationResponseJSON = await startAuthentication({
        optionsJSON: credentialRequestOptions.publicKey,
      });

      const verificationResponse = await fetch(`/login/finish`, {
        method: "POST",
        body: JSON.stringify(assertion),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!verificationResponse.ok) {
        throw new Error(`HTTP error! status: ${verificationResponse.status}`);
      }
      const verificationResponseJSON = await verificationResponse.json();

      if (
        verificationResponseJSON &&
        verificationResponseJSON.status === "ok"
      ) {
        setNotification("Successfully logged in.");
      } else {
        setNotification("Login failed.");
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        switch (err.name) {
          case "TypeError":
            setNotification("An account with this email does not exist");
            break;
          default:
            console.log(err);
            setNotification("Login canceled or Failed due to an error.");
        }
      } else {
        setNotification("Login failed due to an unknown error.");
      }
    }
  }

  return (
    <>
      <div className="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-sm">
          <img
            className="mx-auto h-20 w-auto"
            src="/passkey_logo.png"
            alt="Your Company"
          />
          <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
            Sign in to your account
          </h2>
        </div>
        <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
          <div className="space-y-6">
            <div className="text-sm text-center min-h-5 font-normal text-blue-400">
              {notification}
            </div>
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium leading-6 text-gray-900"
              >
                Email address
              </label>
              <div className="mt-2">
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email webauthn"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>
            </div>
            <div>
              <div className="flex items-center justify-between">
                <label
                  htmlFor="password"
                  className="block text-sm font-medium leading-6 text-gray-900"
                >
                  Password
                </label>
                <div className="text-sm">
                  <a
                    href="#"
                    className="font-semibold text-indigo-600 hover:text-indigo-500"
                  >
                    Forgot password?
                  </a>
                </div>
              </div>
              <div className="mt-2">
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  required
                  className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>
            </div>
            <div>
              <button
                onClick={loginUser}
                type="submit"
                className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                Sign in
              </button>
            </div>
          </div>
          <p className="mt-10 text-center text-sm text-gray-500">
            <a
              href="/sign-up"
              className="font-semibold leading-6 text-indigo-600 hover:text-indigo-500"
            >
              Sign up for a new account
            </a>
          </p>
        </div>
      </div>
    </>
  );
}

export default Login;
