import type { QueryClient } from "@tanstack/react-query";

import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";

import { AuthContextType } from "@/auth";
import { DefaultLayout } from "@/core/layouts";
import { UiProvider } from "@/core/providers";
import ErrorBoundary from "@/core/components/error-boundary";
import NotFound from "@/core/components/404";
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
