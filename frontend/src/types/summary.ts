export interface CategorySummary {
  id: string
  name: string
  contractCount: number
  monthlyTotal: number
}

export interface Summary {
  totalContracts: number
  totalMonthlyAmount: number
  categories: CategorySummary[]
}
