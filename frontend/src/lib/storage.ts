import type { Category } from "@/types/category"
import type { Contract } from "@/types/contract"

const STORAGE_KEY = "contracts-app-data"
const CURRENT_VERSION = 1

interface AppData {
  version: number
  categories: Category[]
  contracts: Contract[]
}

function getDefaultData(): AppData {
  return { version: CURRENT_VERSION, categories: [], contracts: [] }
}

export function readData(): AppData {
  const raw = localStorage.getItem(STORAGE_KEY)
  if (!raw) return getDefaultData()
  try {
    return JSON.parse(raw) as AppData
  } catch {
    return getDefaultData()
  }
}

export function writeData(data: AppData): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
}
