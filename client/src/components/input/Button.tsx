import React from "react";
type Props = {
  onClickFunc: () => void;
  buttonText: string;
};
export function Button({ onClickFunc, buttonText }: Props): React.ReactElement {
  return (
    <button
      onClick={() => onClickFunc()}
      className="flex w-full justify-center rounded-md bg-[#FDDD00] text-black px-3 py-1.5 text-sm font-semibold leading-6 shadow-sm"
    >
      {buttonText}
    </button>
  );
}
