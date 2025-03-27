import React from "react";

type Props = {
  title: string;
  link: string;
};

export function MenuItem({ title, link }: Props): React.ReactElement {
  return (
    <a
      href={link}
      className="flex items-center justify-between w-full text-gray-700 hover:text-indigo-600 transition-colors duration-200 border border-gray-200/100 hover:border-gray-300/50 rounded-lg p-3"
    >
      <span className="text-base font-medium">{title}</span>
      <svg
        className="w-4 h-4 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M9 5l7 7-7 7"
        />
      </svg>
    </a>
  );
}
