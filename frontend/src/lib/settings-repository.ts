import type { Settings } from "@/types/settings"
import { get, put } from "./api"

export async function getSettings(): Promise<Settings> {
  return get<Settings>("/settings")
}

export async function updateSettings(data: Settings): Promise<Settings> {
  return put<Settings>("/settings", data)
}

export async function changePassword(currentPassword: string, newPassword: string): Promise<void> {
  return put<void>("/settings/password", { currentPassword, newPassword })
}
