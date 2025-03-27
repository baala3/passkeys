import React from "react";
import { BackButton } from "../input/BackButton";

type Props = {
  children: React.ReactNode;
  parent?: string;
};
export function Layout({ children, parent }: Props): React.ReactElement {
  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-100 via-purple-50 to-pink-100 flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-md backdrop-blur-xl bg-white/20 rounded-2xl shadow-2xl shadow-indigo-500/20 border border-white/30 p-8">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
            passkeys with go-webauthn
          </h1>
          <div className="h-1 w-20 bg-gradient-to-r from-indigo-600/50 to-purple-600/50 rounded-full mx-auto mt-2"></div>
          <p className="text-gray-600 mt-2">Try it out with the demo below</p>
        </div>
        {parent && <BackButton parent={parent} />}
        {children}
      </div>
    </div>
  );
}
