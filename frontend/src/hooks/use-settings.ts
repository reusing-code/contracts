import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { getSettings, updateSettings, changePassword } from "@/lib/settings-repository"
import type { Settings } from "@/types/settings"

const SETTINGS_KEY = ["settings"] as const

export function useSettings() {
  return useQuery({
    queryKey: SETTINGS_KEY,
    queryFn: getSettings,
  })
}

export function useUpdateSettings() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Settings) => updateSettings(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: SETTINGS_KEY }),
  })
}

export function useChangePassword() {
  return useMutation({
    mutationFn: ({ currentPassword, newPassword }: { currentPassword: string; newPassword: string }) =>
      changePassword(currentPassword, newPassword),
  })
}
