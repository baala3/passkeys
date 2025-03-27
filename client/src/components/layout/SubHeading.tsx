import React from "react";

type Props = {
  children: React.ReactNode;
};

export function SubHeading({ children }: Props): React.ReactElement {
  return (
    <div className="text-center mb-6">
      <h3 className="text-xl font-semibold text-gray-800">{children}</h3>
      <div className="h-0.5 w-16 bg-gradient-to-r from-indigo-500/30 to-purple-500/30 rounded-full mx-auto mt-2"></div>
    </div>
  );
}
