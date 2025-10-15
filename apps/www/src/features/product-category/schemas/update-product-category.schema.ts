import { z } from "zod";

export const updateProductCategorySchema = z.object({
  name: z.string().min(1, "Name is required"),
});

export type UpdateProductCategoryDto = z.infer<
  typeof updateProductCategorySchema
>;
