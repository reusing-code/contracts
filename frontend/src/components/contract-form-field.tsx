import { useTranslation } from "react-i18next"
import type { FieldConfig } from "@/config/contract-fields"
import type { ContractFormData } from "@/types/contract"
import type { Control } from "react-hook-form"
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"

interface ContractFormFieldProps {
  config: FieldConfig
  control: Control<ContractFormData>
}

export function ContractFormField({ config, control }: ContractFormFieldProps) {
  const { t } = useTranslation()

  return (
    <FormField
      control={control}
      name={config.key as keyof ContractFormData}
      render={({ field }) => (
        <FormItem>
          <FormLabel>
            {t(config.i18nKey)}
            {config.required && " *"}
          </FormLabel>
          <FormControl>
            {config.type === "textarea" ? (
              <Textarea
                {...field}
                value={(field.value as string) ?? ""}
                onChange={(e) => field.onChange(e.target.value || undefined)}
              />
            ) : (
              <Input
                type={config.type === "number" ? "number" : config.type === "date" ? "date" : "text"}
                {...field}
                value={field.value === undefined || field.value === null ? "" : String(field.value)}
                onChange={(e) => {
                  if (config.type === "number") {
                    const val = e.target.value
                    field.onChange(val === "" ? undefined : Number(val))
                  } else {
                    field.onChange(e.target.value || undefined)
                  }
                }}
              />
            )}
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  )
}
