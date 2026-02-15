import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { differenceInDays } from "date-fns"
import type { Contract } from "@/types/contract"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getRenewalRowClass(c: Contract): string | undefined {
  if (!c.cancellationDate) return undefined
  const days = differenceInDays(new Date(c.cancellationDate), new Date())
  if (days <= 30) return "bg-destructive/10 hover:bg-destructive/20"
  if (days <= 90) return "bg-yellow-500/10 hover:bg-yellow-500/20"
  return "text-muted-foreground opacity-75"
}
