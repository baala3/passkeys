import React from "react";
import { motion } from "framer-motion";

type Props = {
  onClickFunc: () => void;
  buttonText: string;
};
export function Button({ onClickFunc, buttonText }: Props): React.ReactElement {
  return (
    <motion.button
      whileHover={{ scale: 1.03 }}
      whileTap={{ scale: 0.97 }}
      onClick={() => onClickFunc()}
      className="flex w-full justify-center rounded-md bg-[#FDDD00] text-black px-3 py-1.5 text-sm font-semibold leading-6 shadow-sm cursor-pointer"
    >
      {buttonText}
    </motion.button>
  );
}
