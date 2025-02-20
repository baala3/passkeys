import React from "react";

type Props = {
  onClickFunc: () => void;
  buttonText: string;
};

export function LinkButton({
  onClickFunc,
  buttonText,
}: Props): React.ReactElement {
  return (
    <button
      onClick={() => onClickFunc()}
      className="text-orange-700 px-3 py-1.5 font-semibold leading-6 hover:text-orange-500"
    >
      {buttonText}
    </button>
  );
}
