import React from "react";
import { motion } from "framer-motion";

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
    <motion.div
      layout
      whileHover={{ scale: 1.01 }}
      whileTap={{ scale: 0.99 }}
      transition={{
        duration: 0.8,
        ease: [0, 0.71, 0.2, 1.01],
        layout: {
          duration: 0.8,
          ease: [0, 0.71, 0.2, 1.01],
        },
      }}
    >
      <input
        type={type}
        placeholder={placeholder}
        autoComplete={autoComplete(type)}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-inset focus:ring-black sm:text-sm sm:leading-6"
      />
    </motion.div>
  );
}
