import React, { useState } from "react";
import { RegistrationResponseJSON } from "@simplewebauthn/types";
import { startRegistration } from "@simplewebauthn/browser";

function Register(): React.ReactElement {
  const [email, setEmail] = useState("");
  const [notification, setNotification] = useState("");

  async function registerUser() {
    const isValidEmail = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/g;
    if (!email || !isValidEmail.test(email)) {
      setNotification("Please enter a valid email");
      return;
    }

    try {
      const response = await fetch(`/register/begin`, {
        method: "POST",
        body: JSON.stringify({ email }),
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const credentialCreationOptions = await response.json();
      const registrationResponse: RegistrationResponseJSON =
        await startRegistration({
          optionsJSON: credentialCreationOptions.publicKey,
        });

      const verificationResponse = await fetch(`/register/finish`, {
        method: "POST",
        body: JSON.stringify(registrationResponse),
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
        setNotification("Successfully registered.");
      } else {
        setNotification("Registration failed.");
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        switch (err.name) {
          case "InvalidStateError":
            setNotification("An account with this email already exists");
            break;
          default:
            setNotification(`Registration canceled or Failed due to an error`);
        }
      } else {
        setNotification("Registration failed due to an unknown error.");
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
            Create a new account
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
                  autoComplete="email"
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
                onClick={registerUser}
                className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                Sign up
              </button>
            </div>
          </div>
          <p className="mt-10 text-center text-sm text-gray-500">
            <a
              href="/"
              className="font-semibold leading-6 text-indigo-600 hover:text-indigo-500"
            >
              Already have an account?
            </a>
          </p>
        </div>
      </div>
    </>
  );
}

export default Register;
