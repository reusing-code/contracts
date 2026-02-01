import { createRouter } from "@tanstack/react-router"
import { rootRoute } from "./__root"
import { indexRoute } from "."
import { categoryRoute } from "./categories.$categoryId"
import { loginRoute } from "./login"
import { upcomingRenewalsRoute } from "./upcoming-renewals"
import { settingsRoute } from "./settings"

const routeTree = rootRoute.addChildren([indexRoute, categoryRoute, loginRoute, upcomingRenewalsRoute, settingsRoute])

export const router = createRouter({ routeTree })

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router
  }
}
