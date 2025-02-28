import React from "react";
import { AnimatePresence, motion } from "framer-motion";

type Props = {
  notification: string;
};

export function Notification({ notification }: Props): React.ReactElement {
  return (
    <AnimatePresence mode="wait">
      {notification !== "" && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.3 }}
          className="text-sm text-center font-normal text-blue-400 mb-4"
        >
          {notification}
        </motion.div>
      )}
    </AnimatePresence>
  );
}
