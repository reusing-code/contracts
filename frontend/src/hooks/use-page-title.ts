import { useEffect } from "react"

export function usePageTitle(title: string, appName: string) {
  useEffect(() => {
    document.title = title ? `${title} - ${appName}` : appName
  }, [title, appName])
}
