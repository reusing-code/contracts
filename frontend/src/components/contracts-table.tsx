import { useTranslation } from "react-i18next"
import { format } from "date-fns"
import { ExternalLink, FileText, MoreVertical } from "lucide-react"
import type { Contract } from "@/types/contract"
import { contractFields } from "@/config/contract-fields"
import { getEarliestCancellationDate, isExpired } from "@/lib/contract-utils"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

const tableColumns = contractFields
  .filter((f) => f.showInTable)
  .sort((a, b) => a.tableOrder - b.tableOrder)

interface ContractsTableProps {
  contracts: Contract[]
  onEdit: (contract: Contract) => void
  onDelete: (contract: Contract) => void
}

function formatCellValue(contract: Contract, key: string, currency: string): string {
  const value = contract[key as keyof Contract]
  if (value === undefined || value === null || value === "") return "-"
  if (key === "pricePerMonth") return `${Number(value).toFixed(2)} ${currency}`
  if (key === "startDate" || key === "endDate") return format(new Date(value as string), "yyyy-MM-dd")
  return String(value)
}

export function ContractsTable({ contracts, onEdit, onDelete }: ContractsTableProps) {
  const { t } = useTranslation()
  const currency = t("common.currency")

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            {tableColumns.map((col) => (
              <TableHead key={col.key}>{t(col.i18nKey)}</TableHead>
            ))}
            <TableHead>{t("contract.cancellationDate")}</TableHead>
            <TableHead>{t("contract.links")}</TableHead>
            <TableHead className="w-12" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {contracts.length === 0 ? (
            <TableRow>
              <TableCell colSpan={tableColumns.length + 3} className="text-center text-muted-foreground py-8">
                {t("contract.noContracts")}
              </TableCell>
            </TableRow>
          ) : (
            contracts.map((contract) => {
              const expired = isExpired(contract)
              const cancellation = getEarliestCancellationDate(contract)
              return (
                <TableRow key={contract.id} className={expired ? "opacity-50" : undefined}>
                  {tableColumns.map((col) => (
                    <TableCell key={col.key}>{formatCellValue(contract, col.key, currency)}</TableCell>
                  ))}
                  <TableCell>
                    {expired ? (
                      <Badge variant="secondary">{t("contract.expired")}</Badge>
                    ) : cancellation ? (
                      format(cancellation, "yyyy-MM-dd")
                    ) : (
                      "-"
                    )}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      {contract.customerPortalUrl && (
                        <a href={contract.customerPortalUrl} target="_blank" rel="noopener noreferrer">
                          <Button variant="ghost" size="icon" className="h-8 w-8" title={t("fields.customerPortalUrl")}>
                            <ExternalLink className="h-4 w-4" />
                          </Button>
                        </a>
                      )}
                      {contract.paperlessUrl && (
                        <a href={contract.paperlessUrl} target="_blank" rel="noopener noreferrer">
                          <Button variant="ghost" size="icon" className="h-8 w-8" title={t("fields.paperlessUrl")}>
                            <FileText className="h-4 w-4" />
                          </Button>
                        </a>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="h-8 w-8">
                          <MoreVertical className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => onEdit(contract)}>
                          {t("common.edit")}
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => onDelete(contract)} className="text-destructive">
                          {t("common.delete")}
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              )
            })
          )}
        </TableBody>
      </Table>
    </div>
  )
}
