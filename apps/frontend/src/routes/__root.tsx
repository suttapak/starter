import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";

import { ThemeSwitch } from "@/components/theme-switch";

export const Route = createRootRoute({
  component: () => (
    <>
      <ThemeSwitch />
      <Outlet />
      <TanStackRouterDevtools />
    </>
  ),
});
