import React from "react";

type Props = {
  checked: boolean;
  onChange: () => void;
  label?: string;
};

export function Checkbox({
  checked,
  onChange,
  label,
}: Props): React.ReactElement {
  return (
    <label className="flex items-center gap-2 cursor-pointer group">
      <div className="flex items-center">
        <input
          type="checkbox"
          checked={checked}
          onChange={() => onChange()}
          className="h-4 w-4 rounded-md border-gray-300 text-indigo-600 focus:ring-indigo-500 accent-indigo-600"
        />
      </div>
      {label && (
        <span className="text-sm text-gray-700 group-hover:text-indigo-600 transition-colors duration-200">
          {label}
        </span>
      )}
    </label>
  );
}
