import { Progress } from "@heroui/react";
import { createPortal } from "react-dom";

const MSplashPage = () => {
  return (
    <>
      {createPortal(
        <div className="fixed inset-0 flex justify-center items-center">
          <Progress
            isIndeterminate
            aria-label="Loading..."
            className="max-w-sm"
            size="sm"
          />
        </div>,
        document.body,
      )}
    </>
  );
};

export default MSplashPage;
