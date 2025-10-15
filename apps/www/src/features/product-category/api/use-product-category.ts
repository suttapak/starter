import type { PaginatedResponse, Response } from "@/shared/types";
import type { ProductCategoryResponse } from "../types/product-category";
import type { CreateProductCategoryDto } from "../schemas/create-product-category.schema";
import type { UpdateProductCategoryDto } from "../schemas/update-product-category.schema";

import { useParams, useSearch } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { getJson, postJson, putJson } from "@/core/utils/fetch";

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
