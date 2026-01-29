import { addMonths, isBefore, startOfDay } from "date-fns"
import type { Contract } from "@/types/contract"

export function getEarliestCancellationDate(contract: Contract): Date | null {
  if (contract.endDate) return null

  const start = new Date(contract.startDate)
  const today = startOfDay(new Date())
  const minEnd = addMonths(start, contract.minimumDurationMonths)

  if (contract.extensionDurationMonths === 0) {
    const cancellation = addMonths(minEnd, -contract.noticePeriodMonths)
    return isBefore(cancellation, today) ? today : cancellation
  }

  let periodEnd = minEnd
  while (isBefore(periodEnd, today)) {
    periodEnd = addMonths(periodEnd, contract.extensionDurationMonths)
  }

  const cancellation = addMonths(periodEnd, -contract.noticePeriodMonths)
  if (isBefore(cancellation, today)) {
    periodEnd = addMonths(periodEnd, contract.extensionDurationMonths)
    return addMonths(periodEnd, -contract.noticePeriodMonths)
  }

  return cancellation
}

export function isExpired(contract: Contract): boolean {
  if (!contract.endDate) return false
  return isBefore(new Date(contract.endDate), startOfDay(new Date()))
}
