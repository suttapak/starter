import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authed/admin/_layout/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_authed/admin/_layout/"!</div>
}
