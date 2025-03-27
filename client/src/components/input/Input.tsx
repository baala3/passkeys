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
      return "email";
    case "password":
      return "current-password";
    default:
      return "off";
  }
}

export function Input({
  type,
  placeholder,
  value,
  onChange,
}: Props): React.ReactElement {
  return (
    <div className="relative">
      <input
        type={type}
        placeholder={placeholder}
        autoComplete={autoComplete(type)}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="block w-full rounded-lg border border-gray-200/100 bg-white/50 backdrop-blur-sm p-3 text-gray-700 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500/50 transition-all duration-200 sm:text-sm sm:leading-6"
      />
    </div>
  );
}
