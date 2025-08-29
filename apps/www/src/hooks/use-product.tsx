import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams, useSearch } from "@tanstack/react-router";
import { z } from "zod";
import { t } from "i18next";

import { PaginatedResponse, Response } from "@/types/api_response";
import { deleteJson, getJson, postJson, putJson } from "@/utils/fetch";
import { ProductResponse } from "@/types/product_response";

export const keys = {
  products: (
    id: string,
    page: number,
    limit: number,
    code?: string,
    name?: string,
    uom?: string,
  ) => ["product", id, page, limit, code, name, uom] as const,
  product: (id: string, pId: string) => ["product", id, pId] as const,
  productTransaction: (
    teamId: number,
    productId: number,
    page: number,
    limit: number,
  ) => ["product", "transaction", teamId, productId, page, limit] as const,
};

export const useGetProducts = () => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });
  const { page, limit, code, name, uom } = useSearch({
    from: "/_authed/team/_id/$id/_layout/product/",
  });

  return useQuery({
    queryKey: keys.products(id, page, limit, code, name, uom),
    queryFn: () =>
      getJson<PaginatedResponse<ProductResponse>>(`/teams/${id}/products`, {
        page,
        limit,
        code,
        name,
        uom,
      }),
    enabled: !!id,
  });
};

export const createProductSchema = z.object({
  code: z.string().optional(),
  name: z.string().min(1, "Product name is required"),
  description: z.string().optional(),
  uom: z.string().min(1, "Product uom is required"),
  price: z.number().min(0, "Product price must be a positive number"),
  category_id: z.array(z.number()).optional(),
});

export type CreateProductDto = z.infer<typeof createProductSchema>;

export const useCreateProduct = (
  onSuccess?: (data: ProductResponse) => void,
) => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  const client = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateProductDto) =>
      postJson<Response<ProductResponse>>(`/teams/${id}/products`, data),
    onSuccess: (data) => {
      client.invalidateQueries({ queryKey: ["product", id] });
      onSuccess?.(data.data.data);
    },
  });
};
const MAX_FILE_SIZE = 5 * 1024 * 1024;

// upload files
export const uploadProductImageSchema = z.object({
  product_id: z.number(),
  files: z
    .array(z.instanceof(File))
    .refine((f) => Array.from(f).every((file) => file.size <= MAX_FILE_SIZE)),
});

export type UploadProductImageDto = z.infer<typeof uploadProductImageSchema>;

export const useUploadProductImages = (onSuccess?: () => void) => {
  const client = useQueryClient();
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  return useMutation({
    mutationFn: (data: UploadProductImageDto) => {
      const safe = uploadProductImageSchema.safeParse(data);

      if (!safe.success) throw new Error(t("product.upload_image.error"));
      const form = new FormData();

      data.files?.forEach((file) => {
        form.append("files", file);
      });

      return postJson<Response<void>>(
        `/teams/${id}/products/${data.product_id}/upload_image`,
        form,
      );
    },
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ["product", id] });
      onSuccess?.();
    },
  });
};

export const useDeleteProductImage = (onSuccess?: () => void) => {
  const client = useQueryClient();
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  return useMutation({
    mutationFn: (data: { productId: number; productImageId: number }) =>
      deleteJson<Response<void>>(
        `/teams/${id}/products/${data.productId}/product_image/${data.productImageId}`,
      ),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ["product", id] });
      onSuccess?.();
    },
  });
};

export const createUpdateSchema = z.object({
  code: z.string().min(1, "Product code is required"),
  name: z.string().min(1, "Product name is required"),
  description: z.string().optional(),
  uom: z.string().min(1, "Product uom is required"),
  price: z.number().min(0, "Product price must be a positive number"),
  category_id: z.array(z.number()),
});
export type UpdateProductDto = z.infer<typeof createUpdateSchema>;
export const useUpdateProduct = (pid: number, onSuccess?: () => void) => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  const client = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateProductDto) =>
      putJson(`/teams/${id}/products/${pid}`, data),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ["product", id] });
      onSuccess?.();
    },
  });
};

export const useGetProduct = (pId: string) => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  return useQuery({
    queryKey: keys.product(id, pId),
    queryFn: () =>
      getJson<Response<ProductResponse>>(`/teams/${id}/products/${pId}`),
    enabled: !!pId,
  });
};
