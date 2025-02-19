import React from "react";

type Props = {
  title: string;
  link: string;
};

export function MenuItem({ title, link }: Props): React.ReactElement {
  return (
    <div className="text-xl cursor-pointer p-3 border-b hover:bg-gray-100">
      <a href={link}>{title}</a>
    </div>
  );
}
