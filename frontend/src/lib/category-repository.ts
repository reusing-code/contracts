import type { Category, CategoryFormData } from "@/types/category"
import { readData, writeData } from "./storage"

export async function getAllCategories(): Promise<Category[]> {
  return readData().categories
}

export async function getCategoryById(id: string): Promise<Category | undefined> {
  return readData().categories.find((c) => c.id === id)
}

export async function createCategory(data: CategoryFormData): Promise<Category> {
  const store = readData()
  const now = new Date().toISOString()
  const category: Category = {
    id: crypto.randomUUID(),
    ...data,
    createdAt: now,
    updatedAt: now,
  }
  store.categories.push(category)
  writeData(store)
  return category
}

export async function updateCategory(id: string, data: CategoryFormData): Promise<Category> {
  const store = readData()
  const idx = store.categories.findIndex((c) => c.id === id)
  if (idx === -1) throw new Error(`Category ${id} not found`)
  store.categories[idx] = {
    ...store.categories[idx],
    ...data,
    updatedAt: new Date().toISOString(),
  }
  writeData(store)
  return store.categories[idx]
}

export async function deleteCategory(id: string): Promise<void> {
  const store = readData()
  store.categories = store.categories.filter((c) => c.id !== id)
  store.contracts = store.contracts.filter((c) => c.categoryId !== id)
  writeData(store)
}
