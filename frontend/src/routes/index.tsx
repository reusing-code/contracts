import { createRoute } from "@tanstack/react-router"
import { rootRoute } from "./__root"

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: () => (
    <div className="flex min-h-screen items-center justify-center">
      <h1 className="text-4xl font-bold">Contracts</h1>
    </div>
  ),
})
