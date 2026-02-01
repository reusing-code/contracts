import { useState } from "react"
import { createRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"
import { rootRoute } from "./__root"
import { useUpcomingRenewals } from "@/hooks/use-contracts"
import { useSettings } from "@/hooks/use-settings"
import { updateContract, deleteContract } from "@/lib/contract-repository"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { differenceInDays } from "date-fns"
import type { Contract, ContractFormData } from "@/types/contract"
import { ContractsTable } from "@/components/contracts-table"
import { ContractDialog } from "@/components/contract-dialog"
import { DeleteConfirmDialog } from "@/components/delete-confirm-dialog"

export const upcomingRenewalsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/upcoming-renewals",
  component: UpcomingRenewalsPage,
})

function UpcomingRenewalsPage() {
  const { t } = useTranslation()
  const { data: settings } = useSettings()
  const { data: contracts = [] } = useUpcomingRenewals(settings?.renewalDays)
  const qc = useQueryClient()

  const [editingContract, setEditingContract] = useState<Contract | null>(null)
  const [deletingContract, setDeletingContract] = useState<Contract | null>(null)

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: ContractFormData }) => updateContract(id, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["contracts"] })
      qc.invalidateQueries({ queryKey: ["categories"] })
      toast.success(t("contract.updated"))
      setEditingContract(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: ({ id }: { id: string }) => deleteContract(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["contracts"] })
      qc.invalidateQueries({ queryKey: ["categories"] })
      toast.success(t("contract.deleted"))
      setDeletingContract(null)
    },
  })

  function handleUpdate(data: ContractFormData) {
    if (!editingContract) return
    updateMutation.mutate({ id: editingContract.id, data })
  }

  function handleDelete() {
    if (!deletingContract) return
    deleteMutation.mutate({ id: deletingContract.id })
  }

  function getRowClass(c: Contract) {
    if (!c.cancellationDate) return undefined
    const days = differenceInDays(new Date(c.cancellationDate), new Date())
    
    // Very close/Urgent (<= 30 days) - Red
    if (days <= 30) return "bg-destructive/10 hover:bg-destructive/20"
    
    // Approaching (<= 90 days) - Yellow
    if (days <= 90) return "bg-yellow-500/10 hover:bg-yellow-500/20"
    
    // Far off (> 90 days) - Greyed out / less prominent
    return "text-muted-foreground opacity-75"
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <h1 className="text-2xl font-bold">{t("nav.upcomingRenewals")}</h1>
      </div>

      <ContractsTable
        contracts={contracts}
        onEdit={(c) => setEditingContract(c)}
        onDelete={(c) => setDeletingContract(c)}
        getRowClassName={getRowClass}
      />

      <ContractDialog
        open={!!editingContract}
        onOpenChange={(open) => { if (!open) setEditingContract(null) }}
        contract={editingContract}
        onSubmit={handleUpdate}
      />

      <DeleteConfirmDialog
        open={!!deletingContract}
        onOpenChange={(open) => { if (!open) setDeletingContract(null) }}
        description={t("contract.deleteConfirm", { name: deletingContract?.name ?? "" })}
        onConfirm={handleDelete}
      />
    </div>
  )
}
