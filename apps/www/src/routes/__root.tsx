import type { QueryClient } from "@tanstack/react-query";

import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";

import UiProvider from "@/provider/ui";
import DefaultLayout from "@/layout/default";
import { AuthContextType } from "@/auth";
import ErrorBoundary from "@/components/error-boundary";
import NotFound from "@/components/404";
interface MyRouterContext {
  queryClient: QueryClient;
  auth: AuthContextType;
}

export const Route = createRootRoute<MyRouterContext>({
  component: () => {
    return (
      <UiProvider>
        <DefaultLayout>
          <Outlet />
          <ReactQueryDevtools />
          <TanStackRouterDevtools />
        </DefaultLayout>
      </UiProvider>
    );
  },
  notFoundComponent: NotFound,
  errorComponent: ErrorBoundary,
});
