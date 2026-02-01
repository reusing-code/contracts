import { useState } from "react"
import { createRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"
import { rootRoute } from "./__root"
import { useSettings, useUpdateSettings, useChangePassword } from "@/hooks/use-settings"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"

export const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/settings",
  component: SettingsPage,
})

function SettingsPage() {
  const { t } = useTranslation()
  const { data: settings } = useSettings()
  const updateSettings = useUpdateSettings()
  const changePassword = useChangePassword()

  const [renewalDays, setRenewalDays] = useState<number | null>(null)
  const [currentPassword, setCurrentPassword] = useState("")
  const [newPassword, setNewPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")

  const displayDays = renewalDays ?? settings?.renewalDays ?? 90

  function handleSavePreferences(e: React.FormEvent) {
    e.preventDefault()
    updateSettings.mutate(
      { renewalDays: displayDays },
      {
        onSuccess: () => toast.success(t("settings.saved")),
        onError: () => toast.error(t("settings.saveFailed")),
      },
    )
  }

  function handleChangePassword(e: React.FormEvent) {
    e.preventDefault()
    if (newPassword !== confirmPassword) {
      toast.error(t("auth.passwordMismatch"))
      return
    }
    changePassword.mutate(
      { currentPassword, newPassword },
      {
        onSuccess: () => {
          toast.success(t("settings.passwordChanged"))
          setCurrentPassword("")
          setNewPassword("")
          setConfirmPassword("")
        },
        onError: (err) => {
          toast.error(err.message || t("settings.passwordChangeFailed"))
        },
      },
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">{t("nav.settings")}</h1>

      <Card>
        <CardHeader>
          <CardTitle>{t("settings.preferences")}</CardTitle>
          <CardDescription>{t("settings.preferencesDescription")}</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSavePreferences} className="flex items-end gap-4">
            <div className="space-y-2">
              <Label htmlFor="renewalDays">{t("settings.renewalDays")}</Label>
              <Input
                id="renewalDays"
                type="number"
                min={1}
                max={365}
                value={displayDays}
                onChange={(e) => setRenewalDays(Number(e.target.value))}
                className="w-32"
              />
            </div>
            <Button type="submit" disabled={updateSettings.isPending}>
              {t("common.save")}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t("settings.changePassword")}</CardTitle>
          <CardDescription>{t("settings.changePasswordDescription")}</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleChangePassword} className="space-y-4 max-w-sm">
            <div className="space-y-2">
              <Label htmlFor="currentPassword">{t("settings.currentPassword")}</Label>
              <Input
                id="currentPassword"
                type="password"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="newPassword">{t("settings.newPassword")}</Label>
              <Input
                id="newPassword"
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">{t("auth.confirmPassword")}</Label>
              <Input
                id="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </div>
            <Button type="submit" disabled={changePassword.isPending}>
              {t("settings.changePassword")}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
