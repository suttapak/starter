import React from "react";

import { RegisterForm } from "./ui/register-form";

function Page() {
  return (
    <React.Fragment>
      <div className="flex  w-full items-center justify-center p-6 md:p-10">
        <div className="w-full max-w-sm">
          <RegisterForm />
        </div>
      </div>
    </React.Fragment>
  );
}

export default Page;