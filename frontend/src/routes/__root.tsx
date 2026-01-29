import { createRootRoute, Outlet, Link } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"
import { LanguageSwitcher } from "@/components/language-switcher"
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
    <div className="min-h-screen bg-background">
      <header className="border-b">
        <div className="container mx-auto flex h-14 items-center justify-between px-4">
          <Link to="/" className="text-lg font-semibold">
            {t("app.title")}
          </Link>
          <LanguageSwitcher />
        </div>
      </header>
      <main className="container mx-auto px-4 py-6">
        <Outlet />
      </main>
    </div>
  )
}
