import type { PaginatedResponse, Response } from "@/shared/types";
import type { ProductResponse } from "../types/product";
import type { CreateProductDto } from "../schemas/create-product.schema";
import type { UpdateProductDto } from "../schemas/update-product.schema";

import { t } from "i18next";
import { useParams, useSearch } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  uploadProductImageSchema,
  type UploadProductImageDto,
} from "../schemas/upload-product-image.schema";

import { deleteJson, getJson, postJson, putJson } from "@/core/utils/fetch";

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
