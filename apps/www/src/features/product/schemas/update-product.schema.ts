import { z } from "zod";

export const createUpdateSchema = z.object({
  code: z.string().min(1, "Product code is required"),
  name: z.string().min(1, "Product name is required"),
  description: z.string().optional(),
  uom: z.string().min(1, "Product uom is required"),
  price: z.number().min(0, "Product price must be a positive number"),
  category_id: z.array(z.number()),
});

export type UpdateProductDto = z.infer<typeof createUpdateSchema>;
