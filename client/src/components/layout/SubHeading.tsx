import React from "react";

type Props = {
  children: React.ReactNode;
};

export function SubHeading({ children }: Props): React.ReactElement {
  return <div className="text-xl font-bold mb-4">{children}</div>;
}
