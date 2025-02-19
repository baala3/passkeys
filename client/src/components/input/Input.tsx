import React from "react";
type Props = {
  type: string;
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
};
function autoComplete(type: string): string {
  switch (type) {
    case "email":
      return "email webauthn";
    case "password":
      return "current-password";
    default:
      return "";
  }
}
export function Input({
  type,
  placeholder,
  value,
  onChange,
}: Props): React.ReactElement {
  return (
    <input
      type={type}
      placeholder={placeholder}
      autoComplete={autoComplete(type)}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-inset focus:ring-black sm:text-sm sm:leading-6"
    />
  );
}
