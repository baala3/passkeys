import React from "react";

type Props = {
  children: React.ReactNode;
};

export function Heading({ children }: Props): React.ReactElement {
  return (
    <div className="text-center mb-8">
      <h2 className="text-2xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
        {children}
      </h2>
    </div>
  );
}
