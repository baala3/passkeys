import React from "react";
import { motion } from "framer-motion";

type Props = {
  onClickFunc: () => void;
  buttonText: string;
  className?: string;
};
export function Button({
  onClickFunc,
  buttonText,
  className = "",
}: Props): React.ReactElement {
  return (
    <motion.button
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
      onClick={() => onClickFunc()}
      className={`flex w-full justify-center rounded-lg px-4 py-2.5 text-sm font-medium leading-6 shadow-sm cursor-pointer transition-all duration-200 ease-in-out backdrop-blur-sm bg-indigo-600/10 hover:bg-indigo-600/20 text-indigo-600 border border-indigo-600/20 hover:border-indigo-600/30 hover:shadow-indigo-500/10 ${className}`}
    >
      {buttonText}
    </motion.button>
  );
}
