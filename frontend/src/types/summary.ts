export interface CategorySummary {
  id: string
  name: string
  contractCount: number
  monthlyTotal: number
  yearlyTotal: number
}

export interface Summary {
  totalContracts: number
  totalMonthlyAmount: number
  totalYearlyAmount: number
  categories: CategorySummary[]
}
