import { useState } from "react"
import { createRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"
import { usePageTitle } from "@/hooks/use-page-title"
import { toast } from "sonner"
import { Plus, Upload } from "lucide-react"
import { rootRoute } from "./__root"
import { useCategories, useCreateCategory, useUpdateCategory, useDeleteCategory } from "@/hooks/use-categories"
import { getSummary } from "@/lib/contract-repository"
import { useQuery } from "@tanstack/react-query"
import type { Category, CategoryFormData } from "@/types/category"
import { Button } from "@/components/ui/button"
import { CategoryCard } from "@/components/category-card"
import { CategoryDialog } from "@/components/category-dialog"
import { DeleteConfirmDialog } from "@/components/delete-confirm-dialog"
import { ImportDialog } from "@/components/import-dialog"

export const contractsIndexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/contracts",
  component: ContractsDashboardPage,
})

function ContractsDashboardPage() {
  const { t } = useTranslation()
  usePageTitle(t("nav.contracts"), t("app.title"))
  const { data: categories = [] } = useCategories("contracts")
  const { data: summary } = useQuery({
    queryKey: ["summary"],
    queryFn: getSummary,
  })
  const summaryByCategory = new Map(
    (summary?.categories ?? []).map((s) => [s.id, s]),
  )
  const createCategory = useCreateCategory("contracts")
  const updateCategory = useUpdateCategory("contracts")
  const deleteCategory = useDeleteCategory("contracts")

  const [dialogOpen, setDialogOpen] = useState(false)
  const [importOpen, setImportOpen] = useState(false)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [deletingCategory, setDeletingCategory] = useState<Category | null>(null)

  function handleCreate(data: CategoryFormData) {
    createCategory.mutate(data, { onSuccess: () => toast.success(t("category.created")) })
  }

  function handleUpdate(data: CategoryFormData) {
    if (!editingCategory) return
    updateCategory.mutate(
      { id: editingCategory.id, data },
      { onSuccess: () => toast.success(t("category.updated")) },
    )
    setEditingCategory(null)
  }

  function handleDelete() {
    if (!deletingCategory) return
    deleteCategory.mutate(deletingCategory.id, {
      onSuccess: () => toast.success(t("category.deleted")),
    })
    setDeletingCategory(null)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t("nav.contracts")}</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => setImportOpen(true)}>
            <Upload className="mr-2 h-4 w-4" />
            {t("import.button")}
          </Button>
          <Button onClick={() => setDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            {t("dashboard.newCategory")}
          </Button>
        </div>
      </div>

      {categories.length === 0 ? (
        <p className="text-muted-foreground">{t("dashboard.noCategories")}</p>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {categories.map((cat) => {
            const catSummary = summaryByCategory.get(cat.id)
            return (
              <CategoryCard
                key={cat.id}
                category={cat}
                module="contracts"
                totalAmount={catSummary?.monthlyTotal ?? 0}
                secondaryAmount={catSummary?.yearlyTotal ?? 0}
                itemLabel={t("dashboard.contractCount", { count: catSummary?.contractCount ?? 0 })}
                totalLabel={t("dashboard.monthlyTotal", {
                  amount: `${(catSummary?.monthlyTotal ?? 0).toFixed(2)} ${t("common.currency")}`,
                })}
                secondaryLabel={t("dashboard.yearlyTotal", {
                  amount: `${(catSummary?.yearlyTotal ?? 0).toFixed(2)} ${t("common.currency")}`,
                })}
                onEdit={() => setEditingCategory(cat)}
                onDelete={() => setDeletingCategory(cat)}
              />
            )
          })}
        </div>
      )}

      <CategoryDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        onSubmit={handleCreate}
      />

      <CategoryDialog
        open={!!editingCategory}
        onOpenChange={(open) => { if (!open) setEditingCategory(null) }}
        category={editingCategory}
        onSubmit={handleUpdate}
      />

      <ImportDialog open={importOpen} onOpenChange={setImportOpen} />

      <DeleteConfirmDialog
        open={!!deletingCategory}
        onOpenChange={(open) => { if (!open) setDeletingCategory(null) }}
        description={t("category.deleteConfirm", { name: deletingCategory?.nameKey ? t(deletingCategory.nameKey) : deletingCategory?.name ?? "" })}
        onConfirm={handleDelete}
      />
    </div>
  )
}
