import type { PaginatedResponse } from "@/shared/types";
import type { ReportResponse } from "../types/report";
import type { UploadReportTemplateDto } from "../schemas/upload-report-template.schema";
import type { UpdateReportTemplateDto } from "../schemas/update-report-template.schema";

import { useSearch } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { deleteJson, getJson, postJson, putJson } from "@/core/utils/fetch";

const keys = {
  report: (page: number, limit: number) => ["report", page, limit],
};

export const useFindAllReportTemplate = () => {
  const { page, limit } = useSearch({ from: "/_authed/admin/_layout/report" });

  return useQuery({
    queryKey: keys.report(page, limit),
    queryFn: () =>
      getJson<PaginatedResponse<ReportResponse>>(`/report`, { page, limit }),
  });
};

export const useFindAllReportTemplateNoLimit = () => {
  return useQuery({
    queryKey: keys.report(1, 10000),
    queryFn: () =>
      getJson<PaginatedResponse<ReportResponse>>(`/report`, {
        page: 1,
        limit: 10000,
      }),
  });
};

export const useUploadReportTemplate = (onSuccess?: () => void) => {
  return useMutation({
    mutationFn: (data: UploadReportTemplateDto) => {
      const formFile = new FormData();

      formFile.append("file", data.file);
      formFile.append("name", data.name);

      return postJson("/report", formFile);
    },
    onSuccess: () => {
      onSuccess?.();
    },
  });
};

export const useUpdateReportTemplate = (onSuccess?: () => void) => {
  const { page, limit } = useSearch({ from: "/_authed/admin/_layout/report" });

  const client = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateReportTemplateDto & { templateId: number }) => {
      const formFile = new FormData();

      if (!!data.file) {
        formFile.append("file", data.file);
      }

      formFile.append("name", data.name);

      return putJson(`/report/${data.templateId}`, formFile);
    },
    onSuccess: () => {
      client.invalidateQueries({ queryKey: keys.report(page, limit) });
      onSuccess?.();
    },
  });
};

export const useDeleteReportTemplate = (onSuccess?: () => void) => {
  const { page, limit } = useSearch({ from: "/_authed/admin/_layout/report" });

  const client = useQueryClient();

  return useMutation({
    mutationFn: ({ id }: { id: number }) => deleteJson(`/report/${id}`),
    onSuccess: () => {
      onSuccess?.();
      client.invalidateQueries({ queryKey: keys.report(page, limit) });
    },
  });
};
