import { createRouter } from "@tanstack/react-router"
import { rootRoute } from "./__root"
import { indexRoute } from "."
import { categoryRoute } from "./categories.$categoryId"

const routeTree = rootRoute.addChildren([indexRoute, categoryRoute])

export const router = createRouter({ routeTree })

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router
  }
}
