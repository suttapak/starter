import { createFileRoute, Outlet } from "@tanstack/react-router";

import { NavbarAuthed } from "@/components/navbar-authed";

export const Route = createFileRoute("/_authed/_user")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <>
      <NavbarAuthed />
      <main className="container mx-auto max-w-7xl flex-grow p-2">
        <Outlet />
      </main>
    </>
  );
}
