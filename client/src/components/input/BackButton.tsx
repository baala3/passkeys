import React from "react";

type Props = {
  parent: string;
};

export function BackButton({ parent }: Props): React.ReactElement {
  return (
    <button
      onClick={() => (window.location.href = parent)}
      className="absolute top-4 left-4 text-gray-700 hover:text-indigo-600 transition-all duration-200 flex items-center gap-2 font-semibold hover:scale-150 cursor-pointer"
    >
      <svg
        className="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M10 19l-7-7m0 0l7-7m-7 7h18"
        />
      </svg>
    </button>
  );
}
