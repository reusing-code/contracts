import { createRootRoute, Outlet, Link } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"
import { LanguageSwitcher } from "@/components/language-switcher"
import { Sidebar } from "@/components/sidebar"
import { seedDefaults } from "@/lib/seed"

export const rootRoute = createRootRoute({
  beforeLoad: () => {
    seedDefaults()
  },
  component: RootLayout,
})

function RootLayout() {
  const { t } = useTranslation()

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="border-b">
        <div className="flex h-14 items-center justify-between px-4">
          <Link to="/" className="text-lg font-semibold">
            {t("app.title")}
          </Link>
          <LanguageSwitcher />
        </div>
      </header>
      <div className="flex flex-1">
        <Sidebar />
        <main className="flex-1 px-6 py-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
