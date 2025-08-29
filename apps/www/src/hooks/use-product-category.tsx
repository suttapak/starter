import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams, useSearch } from "@tanstack/react-router";
import { z } from "zod";

import { getJson, postJson, putJson } from "@/utils/fetch";
import { PaginatedResponse, Response } from "@/types/api_response";
import { ProductCategoryResponse } from "@/types/product_category_response";

const keys = {
  category: (id: string, page: number, limit: number) =>
    ["product-category", id, page, limit] as const,
  categoryInProduct: (id: string) =>
    ["product-category-in-product", id] as const,
};

export const useGetProductCategory = () => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });
  const { page = 1, limit = 10 } = useSearch({
    from: "/_authed/team/_id/$id/_layout/product/category",
  });

  return useQuery({
    queryKey: keys.category(id, page, limit),
    queryFn: () =>
      getJson<PaginatedResponse<ProductCategoryResponse>>(
        `/teams/${id}/product_category`,
        { page, limit },
      ),
    enabled: !!id,
  });
};

export const useGetProductCategoryInProductPage = () => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  return useQuery({
    queryKey: keys.categoryInProduct(id),
    queryFn: () =>
      getJson<PaginatedResponse<ProductCategoryResponse>>(
        `/teams/${id}/product_category`,
        { page: 1, limit: 1000 },
      ),
    enabled: !!id,
  });
};

export const createProductCategorySchema = z.object({
  name: z.string().min(1, "Name is required"),
});

export type CreateProductCategoryDto = z.infer<
  typeof createProductCategorySchema
>;

export const useCreateProductCategory = (onSuccess?: () => void) => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateProductCategoryDto) =>
      postJson<Response<never>>(`/teams/${id}/product_category`, {
        ...data,
        team_id: Number(id),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-category", id] });
      onSuccess?.();
    },
  });
};

export const updateProductCategorySchema = z.object({
  name: z.string().min(1, "Name is required"),
});

export type UpdateProductCategoryDto = z.infer<
  typeof updateProductCategorySchema
>;

export const useUpdateProductCategory = (
  categoryId: number,
  onSuccess?: () => void,
) => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateProductCategoryDto) =>
      putJson<UpdateProductCategoryDto, Response<never>>(
        `/teams/${id}/product_category/${categoryId}`,
        data,
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-category", id] });
      onSuccess?.();
    },
  });
};
