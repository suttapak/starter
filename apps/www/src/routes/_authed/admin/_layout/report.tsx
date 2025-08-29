import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import { Pagination } from "@heroui/react";

import { useFindAllReportTemplate } from "@/hooks/use-report";
import MSplashPage from "@/components/splash";
import UploadReportTemplateModal from "@/components/upload-report-template-modal";
import ReportTemplateCard from "@/components/report-template-card";
const validateSearch = z.object({
  page: z.number().prefault(1),
  limit: z.number().prefault(10),
});

export const Route = createFileRoute("/_authed/admin/_layout/report")({
  component: RouteComponent,
  validateSearch,
});

function RouteComponent() {
  const { page, limit } = Route.useSearch();
  const navigate = Route.useNavigate();
  const { data, isLoading } = useFindAllReportTemplate();

  if (isLoading) return <MSplashPage />;

  return (
    <div className="flex flex-col gap-2">
      <div className="flex justify-end">
        <UploadReportTemplateModal />
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-4 gap-4">
        {data?.data.data.map((item) => (
          <ReportTemplateCard key={item.id} item={item} />
        ))}
      </div>
      <div className="flex justify-center items-center gap-2">
        <Pagination
          isCompact
          showControls
          showShadow
          color="secondary"
          page={page}
          total={data?.data.meta.total_page || 0}
          onChange={(page) => navigate({ search: { page: page } })}
        />
        <label className="flex items-center text-default-400 text-small">
          Rows per page:
          <select
            className="bg-transparent outline-none text-default-400 text-small"
            defaultValue={limit}
            onChange={(e) => {
              const newLimit = parseInt(e.target.value, 10);

              navigate({ search: { limit: newLimit, page: 1 } });
            }}
          >
            <option value="5">5</option>
            <option value="10">10</option>
            <option value="20">20</option>
            <option value="50">50</option>
          </select>
        </label>
      </div>
    </div>
  );
}
