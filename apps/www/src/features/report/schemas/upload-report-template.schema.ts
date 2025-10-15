import { z } from "zod";

export const uploadReportTemplateSchema = z.object({
  file: z
    .instanceof(File)
    .refine((file) => file && file.name.toLowerCase().endsWith(".odt")),
  name: z.string(),
});

export type UploadReportTemplateDto = z.infer<
  typeof uploadReportTemplateSchema
>;
