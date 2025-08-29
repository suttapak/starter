import { createFileRoute, Outlet } from "@tanstack/react-router";

import { NavbarPublic } from "@/components/navbar-public";

export const Route = createFileRoute("/_public")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <>
      <NavbarPublic />
      <div className="container mx-auto p-2">
        <Outlet />
      </div>
    </>
  );
}
