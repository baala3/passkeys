import React from "react";
export function Divider(): React.ReactElement {
  return (
    <div className="m-6 py-3 flex items-center text-xs text-gray-400 uppercase before:flex-1 before:border-t before:border-gray-200 before:me-6 after:flex-1 after:border-t after:border-gray-200 after:ms-6">
      Or
    </div>
  );
}
