import type { Category, CategoryFormData } from "@/types/category"
import { del, get, post, put } from "./api"

export async function getAllCategories(module: string): Promise<Category[]> {
  return get<Category[]>(`/modules/${module}/categories`)
}

export async function getCategoryById(module: string, id: string): Promise<Category> {
  return get<Category>(`/modules/${module}/categories/${id}`)
}

export async function createCategory(module: string, data: CategoryFormData): Promise<Category> {
  return post<Category>(`/modules/${module}/categories`, data)
}

export async function updateCategory(module: string, id: string, data: CategoryFormData): Promise<Category> {
  return put<Category>(`/modules/${module}/categories/${id}`, data)
}

export async function deleteCategory(module: string, id: string): Promise<void> {
  return del(`/modules/${module}/categories/${id}`)
}
