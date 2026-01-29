import { z } from "zod/v4"

export const contractSchema = z.object({
  id: z.string().uuid(),
  categoryId: z.string().uuid(),
  name: z.string().min(1),
  productName: z.string().optional(),
  company: z.string().optional(),
  contractNumber: z.string().optional(),
  customerNumber: z.string().optional(),
  pricePerMonth: z.number().nonnegative().optional(),
  startDate: z.string().date(),
  endDate: z.string().date().optional(),
  minimumDurationMonths: z.number().int().nonnegative(),
  extensionDurationMonths: z.number().int().nonnegative(),
  noticePeriodMonths: z.number().int().nonnegative(),
  customerPortalUrl: z.string().url().optional().or(z.literal("")),
  paperlessUrl: z.string().url().optional().or(z.literal("")),
  comments: z.string().optional(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
})

export type Contract = z.infer<typeof contractSchema>

export const contractFormSchema = z.object({
  name: z.string().min(1),
  productName: z.string().optional(),
  company: z.string().optional(),
  contractNumber: z.string().optional(),
  customerNumber: z.string().optional(),
  pricePerMonth: z.number().nonnegative().optional(),
  startDate: z.string().date(),
  endDate: z.string().date().optional(),
  minimumDurationMonths: z.number().int().nonnegative(),
  extensionDurationMonths: z.number().int().nonnegative(),
  noticePeriodMonths: z.number().int().nonnegative(),
  customerPortalUrl: z.string().url().optional().or(z.literal("")),
  paperlessUrl: z.string().url().optional().or(z.literal("")),
  comments: z.string().optional(),
})

export type ContractFormData = z.infer<typeof contractFormSchema>
