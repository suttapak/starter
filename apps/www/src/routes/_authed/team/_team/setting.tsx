import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_authed/team/_team/setting")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_authed/team/_team/setting"!</div>;
}
