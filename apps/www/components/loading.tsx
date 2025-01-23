import { Progress } from "@heroui/react";
import React from "react";
import { createPortal } from "react-dom";

function Loading() {
  return (
    <React.Fragment>
      {createPortal(
        <>
          <div className="fixed inset-0 flex justify-center items-center">
            <Progress
              isIndeterminate
              aria-label="Loading..."
              className="max-w-sm"
              size="sm"
            />
          </div>
        </>,
        document.body,
      )}
    </React.Fragment>
  );
}

export default Loading;
