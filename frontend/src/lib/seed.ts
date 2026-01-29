import { readData, writeData } from "./storage"

const DEFAULT_CATEGORIES = [
  { name: "Insurance", nameKey: "categoryNames.insurance" },
  { name: "Banking / Portfolios", nameKey: "categoryNames.banking" },
  { name: "Telecommunications", nameKey: "categoryNames.telecommunications" },
]

export function seedDefaults(): void {
  const store = readData()
  if (store.categories.length > 0) return

  const now = new Date().toISOString()
  store.categories = DEFAULT_CATEGORIES.map(({ name, nameKey }) => ({
    id: crypto.randomUUID(),
    name,
    nameKey,
    createdAt: now,
    updatedAt: now,
  }))
  writeData(store)
}
