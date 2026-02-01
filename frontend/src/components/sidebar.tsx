import { Link, useMatchRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"
import { useCategories } from "@/hooks/use-categories"
import { cn } from "@/lib/utils"

export function Sidebar() {
  const { t } = useTranslation()
  const { data: categories = [] } = useCategories()
  const matchRoute = useMatchRoute()

  const isDashboard = matchRoute({ to: "/" })

  return (
    <aside className="w-64 shrink-0 border-r bg-background">
      <nav className="flex flex-col gap-1 p-4">
        <Link
          to="/"
          className={cn(
            "rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent",
            isDashboard && "bg-accent",
          )}
        >
          {t("nav.dashboard")}
        </Link>

        <Link
          to="/upcoming-renewals"
          className={cn(
            "rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent",
            matchRoute({ to: "/upcoming-renewals" }) && "bg-accent",
          )}
        >
          {t("nav.upcomingRenewals")}
        </Link>

        <Link
          to="/settings"
          className={cn(
            "rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent",
            matchRoute({ to: "/settings" }) && "bg-accent",
          )}
        >
          {t("nav.settings")}
        </Link>

        <h2 className="mt-4 mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
          {t("nav.categories")}
        </h2>

        {categories.map((category) => {
          const isActive = matchRoute({
            to: "/categories/$categoryId",
            params: { categoryId: category.id },
          })
          return (
            <Link
              key={category.id}
              to="/categories/$categoryId"
              params={{ categoryId: category.id }}
              className={cn(
                "rounded-md px-3 py-2 text-sm transition-colors hover:bg-accent",
                isActive && "bg-accent font-medium",
              )}
            >
              {category.nameKey ? t(category.nameKey) : category.name}
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}
