import { z } from "zod";

export const createProductSchema = z.object({
  code: z.string().optional(),
  name: z.string().min(1, "Product name is required"),
  description: z.string().optional(),
  uom: z.string().min(1, "Product uom is required"),
  price: z.number().min(0, "Product price must be a positive number"),
  category_id: z.array(z.number()).optional(),
});

export type CreateProductDto = z.infer<typeof createProductSchema>;
