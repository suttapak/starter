import { z } from "zod";

export const updateReportTemplateSchema = z.object({
  file: z
    .instanceof(File)
    .optional()
    .refine((file) => !file || file.name.toLowerCase().endsWith(".odt")),
  name: z.string(),
});

export type UpdateReportTemplateDto = z.infer<
  typeof updateReportTemplateSchema
>;
