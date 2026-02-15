import { createRouter } from "@tanstack/react-router"
import { rootRoute } from "./__root"
import { indexRoute } from "."
import { contractsIndexRoute } from "./contracts.index"
import { contractsCategoryRoute } from "./contracts.categories.$categoryId"
import { contractsUpcomingRenewalsRoute } from "./contracts.upcoming-renewals"
import { purchasesIndexRoute } from "./purchases.index"
import { purchasesCategoryRoute } from "./purchases.categories.$categoryId"
import { loginRoute } from "./login"
import { settingsRoute } from "./settings"

const routeTree = rootRoute.addChildren([
  indexRoute,
  contractsIndexRoute,
  contractsCategoryRoute,
  contractsUpcomingRenewalsRoute,
  purchasesIndexRoute,
  purchasesCategoryRoute,
  loginRoute,
  settingsRoute,
])

export const router = createRouter({ routeTree })

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router
  }
}
