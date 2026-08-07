import { FileText } from "lucide-react"
import { usePaperlessThumbnail } from "@/hooks/use-paperless"
import { cn } from "@/lib/utils"

export function PaperlessThumbnail({ documentId, className }: { documentId: number; className?: string }) {
  const url = usePaperlessThumbnail(documentId, true)

  if (!url) {
    return (
      <div className={cn("flex items-center justify-center bg-muted animate-pulse rounded", className)}>
        <FileText className="h-5 w-5 text-muted-foreground" />
      </div>
    )
  }
  return <img src={url} alt="" className={cn("object-cover rounded border", className)} />
}
