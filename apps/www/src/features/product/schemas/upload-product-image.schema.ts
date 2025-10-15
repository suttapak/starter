import { z } from "zod";

const MAX_FILE_SIZE = 5 * 1024 * 1024;

export const uploadProductImageSchema = z.object({
  product_id: z.number(),
  files: z
    .array(z.instanceof(File))
    .refine((f) => Array.from(f).every((file) => file.size <= MAX_FILE_SIZE)),
});

export type UploadProductImageDto = z.infer<typeof uploadProductImageSchema>;
