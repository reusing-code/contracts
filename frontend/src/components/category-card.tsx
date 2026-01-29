import { useTranslation } from "react-i18next"
import { useNavigate } from "@tanstack/react-router"
import { MoreVertical } from "lucide-react"
import type { Category } from "@/types/category"
import type { Contract } from "@/types/contract"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

interface CategoryCardProps {
  category: Category
  contracts: Contract[]
  onEdit: () => void
  onDelete: () => void
}

export function CategoryCard({ category, contracts, onEdit, onDelete }: CategoryCardProps) {
  const { t } = useTranslation()
  const navigate = useNavigate()

  const monthlyTotal = contracts.reduce((sum, c) => sum + (c.pricePerMonth ?? 0), 0)

  return (
    <Card
      className="cursor-pointer transition-colors hover:bg-accent/50"
      onClick={() => navigate({ to: "/categories/$categoryId", params: { categoryId: category.id } })}
    >
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-lg font-semibold">
          {category.nameKey ? t(category.nameKey) : category.name}
        </CardTitle>
        <DropdownMenu>
          <DropdownMenuTrigger asChild onClick={(e) => e.stopPropagation()}>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" onClick={(e) => e.stopPropagation()}>
            <DropdownMenuItem onClick={onEdit}>{t("common.edit")}</DropdownMenuItem>
            <DropdownMenuItem onClick={onDelete} className="text-destructive">
              {t("common.delete")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground">
          {t("dashboard.contractCount", { count: contracts.length })}
        </p>
        {monthlyTotal > 0 && (
          <p className="text-sm font-medium mt-1">
            {t("dashboard.monthlyTotal", {
              amount: `${monthlyTotal.toFixed(2)} ${t("common.currency")}`,
            })}
          </p>
        )}
      </CardContent>
    </Card>
  )
}
