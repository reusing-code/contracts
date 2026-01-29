import type { Contract, ContractFormData } from "@/types/contract"
import { readData, writeData } from "./storage"

export async function getAllContracts(): Promise<Contract[]> {
  return readData().contracts
}

export async function getContractsByCategory(categoryId: string): Promise<Contract[]> {
  return readData().contracts.filter((c) => c.categoryId === categoryId)
}

export async function getContractById(id: string): Promise<Contract | undefined> {
  return readData().contracts.find((c) => c.id === id)
}

export async function createContract(categoryId: string, data: ContractFormData): Promise<Contract> {
  const store = readData()
  const now = new Date().toISOString()
  const contract: Contract = {
    id: crypto.randomUUID(),
    categoryId,
    ...data,
    createdAt: now,
    updatedAt: now,
  }
  store.contracts.push(contract)
  writeData(store)
  return contract
}

export async function updateContract(id: string, data: ContractFormData): Promise<Contract> {
  const store = readData()
  const idx = store.contracts.findIndex((c) => c.id === id)
  if (idx === -1) throw new Error(`Contract ${id} not found`)
  store.contracts[idx] = {
    ...store.contracts[idx],
    ...data,
    updatedAt: new Date().toISOString(),
  }
  writeData(store)
  return store.contracts[idx]
}

export async function deleteContract(id: string): Promise<void> {
  const store = readData()
  store.contracts = store.contracts.filter((c) => c.id !== id)
  writeData(store)
}
