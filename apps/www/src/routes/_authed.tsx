import {
  createFileRoute,
  Navigate,
  Outlet,
  useRouter,
} from "@tanstack/react-router";

import { useAuth } from "@/auth";
import MSplashPage from "@/components/splash";

export const Route = createFileRoute("/_authed")({
  component: RouteComponent,
});

function RouteComponent() {
  const { isInit, isAuthenticated } = useAuth();
  const { state } = useRouter();

  if (!isInit) return <MSplashPage />;
  if (!isAuthenticated)
    return (
      <Navigate
        replace
        search={{ redirect: state.location.href }}
        to="/login"
      />
    );

  return <Outlet />;
}
