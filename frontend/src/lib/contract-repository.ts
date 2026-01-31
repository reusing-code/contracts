import type { Contract, ContractFormData } from "@/types/contract"
import type { Summary } from "@/types/summary"
import { del, get, post, put } from "./api"

export async function getAllContracts(): Promise<Contract[]> {
  return get<Contract[]>("/contracts")
}

export async function getContractsByCategory(categoryId: string): Promise<Contract[]> {
  return get<Contract[]>(`/categories/${categoryId}/contracts`)
}

export async function getContractById(id: string): Promise<Contract | undefined> {
  return get<Contract>(`/contracts/${id}`)
}

export async function createContract(categoryId: string, data: ContractFormData): Promise<Contract> {
  return post<Contract>(`/categories/${categoryId}/contracts`, data)
}

export async function updateContract(id: string, data: ContractFormData): Promise<Contract> {
  return put<Contract>(`/contracts/${id}`, data)
}

export async function deleteContract(id: string): Promise<void> {
  return del(`/contracts/${id}`)
}

export async function getSummary(): Promise<Summary> {
  return get<Summary>("/summary")
}
