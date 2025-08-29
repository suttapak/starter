import { createFileRoute, Navigate } from "@tanstack/react-router";

export const Route = createFileRoute("/_authed/team/_id/$id/_layout/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = Route.useParams();

  return <Navigate replace params={{ id: id }} to="/team/$id/product" />;
}
