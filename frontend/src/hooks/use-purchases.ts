import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import {
  getPurchasesByCategory,
  createPurchase,
  updatePurchase,
  deletePurchase,
} from "@/lib/purchase-repository"
import type { PurchaseFormData } from "@/types/purchase"

const purchasesKey = (categoryId: string) => ["purchases", categoryId] as const

export function useCategoryPurchases(categoryId: string) {
  return useQuery({
    queryKey: purchasesKey(categoryId),
    queryFn: () => getPurchasesByCategory(categoryId),
  })
}

export function useCreatePurchase(categoryId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: PurchaseFormData) => createPurchase(categoryId, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: purchasesKey(categoryId) })
      qc.invalidateQueries({ queryKey: ["categories", "purchases"] })
      qc.invalidateQueries({ queryKey: ["purchases-summary"] })
    },
  })
}

export function useUpdatePurchase(categoryId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: PurchaseFormData }) => updatePurchase(id, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: purchasesKey(categoryId) })
      qc.invalidateQueries({ queryKey: ["categories", "purchases"] })
      qc.invalidateQueries({ queryKey: ["purchases-summary"] })
    },
  })
}

export function useDeletePurchase(categoryId: string) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deletePurchase(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: purchasesKey(categoryId) })
      qc.invalidateQueries({ queryKey: ["categories", "purchases"] })
      qc.invalidateQueries({ queryKey: ["purchases-summary"] })
    },
  })
}
