import { useState } from "react"
import { Link, useMatchRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"
import { useCategories } from "@/hooks/use-categories"
import { cn } from "@/lib/utils"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { ChevronRight, Home } from "lucide-react"

function SidebarSection({
  title,
  defaultOpen = true,
  children,
}: {
  title: string
  defaultOpen?: boolean
  children: React.ReactNode
}) {
  const [open, setOpen] = useState(defaultOpen)

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <CollapsibleTrigger className="flex w-full items-center gap-1 px-3 py-1.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors">
        <ChevronRight className={cn("h-3 w-3 transition-transform", open && "rotate-90")} />
        {title}
      </CollapsibleTrigger>
      <CollapsibleContent>
        <div className="flex flex-col gap-1">
          {children}
        </div>
      </CollapsibleContent>
    </Collapsible>
  )
}

export function Sidebar() {
  const { t } = useTranslation()
  const { data: contractCategories = [] } = useCategories("contracts")
  const { data: purchaseCategories = [] } = useCategories("purchases")
  const matchRoute = useMatchRoute()

  return (
    <aside className="w-64 shrink-0 border-r bg-background">
      <nav className="flex flex-col gap-1 p-4">
        <Link
          to="/"
          className={cn(
            "flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent",
            matchRoute({ to: "/" }) && "bg-accent",
          )}
        >
          <Home className="h-4 w-4" />
          {t("home.title")}
        </Link>

        <div className="my-2" />

        <SidebarSection title={t("nav.contracts")}>
          <Link
            to="/contracts"
            className={cn(
              "rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent",
              matchRoute({ to: "/contracts" }) && "bg-accent",
            )}
          >
            {t("nav.dashboard")}
          </Link>

          <Link
            to="/contracts/upcoming-renewals"
            className={cn(
              "rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent",
              matchRoute({ to: "/contracts/upcoming-renewals" }) && "bg-accent",
            )}
          >
            {t("nav.upcomingRenewals")}
          </Link>

          {contractCategories.map((category) => {
            const isActive = matchRoute({
              to: "/contracts/categories/$categoryId",
              params: { categoryId: category.id },
            })
            return (
              <Link
                key={category.id}
                to="/contracts/categories/$categoryId"
                params={{ categoryId: category.id }}
                className={cn(
                  "rounded-md px-3 py-2 text-sm transition-colors hover:bg-accent pl-6",
                  isActive && "bg-accent font-medium",
                )}
              >
                {category.nameKey ? t(category.nameKey) : category.name}
              </Link>
            )
          })}
        </SidebarSection>

        <SidebarSection title={t("nav.purchases")}>
          <Link
            to="/purchases"
            className={cn(
              "rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent",
              matchRoute({ to: "/purchases" }) && "bg-accent",
            )}
          >
            {t("nav.dashboard")}
          </Link>

          {purchaseCategories.map((category) => {
            const isActive = matchRoute({
              to: "/purchases/categories/$categoryId",
              params: { categoryId: category.id },
            })
            return (
              <Link
                key={category.id}
                to="/purchases/categories/$categoryId"
                params={{ categoryId: category.id }}
                className={cn(
                  "rounded-md px-3 py-2 text-sm transition-colors hover:bg-accent pl-6",
                  isActive && "bg-accent font-medium",
                )}
              >
                {category.nameKey ? t(category.nameKey) : category.name}
              </Link>
            )
          })}
        </SidebarSection>

        <SidebarSection title={t("nav.general")}>
          <Link
            to="/settings"
            className={cn(
              "rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent",
              matchRoute({ to: "/settings" }) && "bg-accent",
            )}
          >
            {t("nav.settings")}
          </Link>
        </SidebarSection>
      </nav>
    </aside>
  )
}
