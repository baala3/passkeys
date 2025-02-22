import React from "react";
import gopher from "../../assets/gopher.png";
type Props = {
  children: React.ReactNode;
};
export function Layout({ children }: Props): React.ReactElement {
  return (
    <div className="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8 text-center font-normal mx-auto">
      <div
        className="w-full bg-[#027D9C] text-white p-4 fixed top-0 left-0 text-xl text-center cursor-pointer"
        onClick={() => {
          window.location.href = "/";
        }}
      >
        Passkeys with go-webauthn
      </div>
      <div className="mt-20 sm:mx-auto sm:w-full sm:max-w-sm">{children}</div>
      <img
        src={gopher}
        alt="gopher"
        title="Go Gopher by Renée French"
        className="fixed bottom-0 right-0 h-24 object-scale-down"
      />
    </div>
  );
}
