import React from "react";
type Props = {
  href: string;
  linkText: string;
};
export function Link({ href, linkText }: Props): React.ReactElement {
  return (
    <p className="m-2 text-sm">
      <a
        href={href}
        className="font-semibold leading-6 text-indigo-600 hover:text-indigo-500 hover:underline"
      >
        {linkText}
      </a>
    </p>
  );
}
