import { useState } from "react"
import { createRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"
import { Plus } from "lucide-react"
import { rootRoute } from "./__root"
import { useCategories, useCreateCategory, useUpdateCategory, useDeleteCategory } from "@/hooks/use-categories"
import { getAllContracts } from "@/lib/contract-repository"
import { useQuery } from "@tanstack/react-query"
import type { Category, CategoryFormData } from "@/types/category"
import { Button } from "@/components/ui/button"
import { CategoryCard } from "@/components/category-card"
import { CategoryDialog } from "@/components/category-dialog"
import { DeleteConfirmDialog } from "@/components/delete-confirm-dialog"

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: DashboardPage,
})

function DashboardPage() {
  const { t } = useTranslation()
  const { data: categories = [] } = useCategories()
  const { data: allContracts = [] } = useQuery({
    queryKey: ["contracts-all"],
    queryFn: getAllContracts,
  })
  const createCategory = useCreateCategory()
  const updateCategory = useUpdateCategory()
  const deleteCategory = useDeleteCategory()

  const [dialogOpen, setDialogOpen] = useState(false)
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
        <h1 className="text-2xl font-bold">{t("dashboard.title")}</h1>
        <Button onClick={() => setDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          {t("dashboard.newCategory")}
        </Button>
      </div>

      {categories.length === 0 ? (
        <p className="text-muted-foreground">{t("dashboard.noCategories")}</p>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {categories.map((cat) => (
            <CategoryCard
              key={cat.id}
              category={cat}
              contracts={allContracts.filter((c) => c.categoryId === cat.id)}
              onEdit={() => {
                setEditingCategory(cat)
              }}
              onDelete={() => {
                setDeletingCategory(cat)
              }}
            />
          ))}
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

      <DeleteConfirmDialog
        open={!!deletingCategory}
        onOpenChange={(open) => { if (!open) setDeletingCategory(null) }}
        description={t("category.deleteConfirm", { name: deletingCategory?.nameKey ? t(deletingCategory.nameKey) : deletingCategory?.name ?? "" })}
        onConfirm={handleDelete}
      />
    </div>
  )
}
