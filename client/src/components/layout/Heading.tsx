import React from "react";

type Props = {
  children: React.ReactNode;
};

export function Heading({ children }: Props): React.ReactElement {
  return <div className="text-4xl font-bold mb-8">{children}</div>;
}
