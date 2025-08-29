import { createFileRoute } from "@tanstack/react-router";
import {
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  Pagination,
  Spinner,
} from "@heroui/react";
import { z } from "zod";
import { useTranslation } from "react-i18next";

import { useGetProductCategory } from "@/hooks/use-product-category";
import CategoryEditModal from "@/components/category-edit-modal";
import CategoryCreateModal from "@/components/category-create-modal";

const categorySearchParamsSchema = z.object({
  page: z.number().prefault(1),
  limit: z.number().prefault(10),
});

export const Route = createFileRoute(
  "/_authed/team/_id/$id/_layout/product/category",
)({
  component: RouteComponent,
  validateSearch: categorySearchParamsSchema,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { page, limit } = Route.useSearch();
  const { data, isLoading } = useGetProductCategory();
  const navigate = Route.useNavigate();

  return (
    <div className="flex flex-col gap-2">
      <div className="flex justify-end gap-2">
        <CategoryCreateModal />
      </div>
      <Table
        bottomContent={
          data?.data.meta && (
            <div className="flex w-full justify-center gap-2 flex-col sm:flex-row">
              <Pagination
                isCompact
                showControls
                showShadow
                color="secondary"
                page={page || 1}
                total={data?.data.meta.total_page}
                onChange={(page) => navigate({ search: { page } })}
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
                  <option value="20">15</option>
                  <option value="50">15</option>
                </select>
              </label>
            </div>
          )
        }
      >
        <TableHeader>
          <TableColumn width={80}>{t("category.state.index")}</TableColumn>
          <TableColumn>{t("category.state.name")}</TableColumn>
          <TableColumn width={80}>Action</TableColumn>
        </TableHeader>
        <TableBody
          items={data?.data.data || []}
          loadingContent={<Spinner />}
          loadingState={isLoading ? "loading" : "idle"}
        >
          {(item) => (
            <TableRow key={item.id}>
              <TableCell>{item.id}</TableCell>
              <TableCell>{item.name}</TableCell>
              <TableCell>
                <CategoryEditModal item={item} />
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
