import React from "react";
import { motion } from "framer-motion";

type Props = {
  onClickFunc: () => void;
  buttonText: string;
};

export function LinkButton({
  onClickFunc,
  buttonText,
}: Props): React.ReactElement {
  return (
    <motion.button
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      onClick={() => onClickFunc()}
      className="text-orange-700 px-3 py-1.5 font-semibold leading-6 hover:text-orange-500"
    >
      {buttonText}
    </motion.button>
  );
}
