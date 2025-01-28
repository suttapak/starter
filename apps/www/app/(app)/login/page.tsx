import React from "react";

import LoginForm from "./ui/login-form";

const Page = () => {
  return (
    <div className="flex  w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <LoginForm />
      </div>
    </div>
  );
};

export default Page;
