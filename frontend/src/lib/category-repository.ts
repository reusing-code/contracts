import type { Category, CategoryFormData } from "@/types/category"
import { del, get, post, put } from "./api"

export async function getAllCategories(): Promise<Category[]> {
  return get<Category[]>("/categories")
}

export async function getCategoryById(id: string): Promise<Category | undefined> {
  return get<Category>(`/categories/${id}`)
}

export async function createCategory(data: CategoryFormData): Promise<Category> {
  return post<Category>("/categories", data)
}

export async function updateCategory(id: string, data: CategoryFormData): Promise<Category> {
  return put<Category>(`/categories/${id}`, data)
}

export async function deleteCategory(id: string): Promise<void> {
  return del(`/categories/${id}`)
}
