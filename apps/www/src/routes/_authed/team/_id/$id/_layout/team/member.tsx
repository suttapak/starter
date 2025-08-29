import { createFileRoute } from "@tanstack/react-router";
import {
  Pagination,
  TableHeader,
  TableColumn,
  TableBody,
  Spinner,
  TableRow,
  TableCell,
  Table,
  Button,
} from "@heroui/react";
import { z } from "zod";
import { MoreHorizontal } from "lucide-react";
import { useTranslation } from "react-i18next";

import { useGetTeamMembers } from "@/hooks/use-team";

const validateSearch = z.object({
  page: z.number().prefault(1),
  limit: z.number().prefault(10),
  code: z.string().optional(),
  name: z.string().optional(),
  uom: z.string().optional(),
});

export const Route = createFileRoute(
  "/_authed/team/_id/$id/_layout/team/member",
)({
  component: RouteComponent,
  validateSearch,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { data, isLoading } = useGetTeamMembers();
  const { page, limit } = Route.useSearch();
  const navigate = Route.useNavigate();

  return (
    <div>
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
          <TableColumn width={80}>{t("team.member.state.index")}</TableColumn>
          <TableColumn width={200}>{t("team.member.state.name")}</TableColumn>
          <TableColumn width={200}>{t("team.member.state.email")}</TableColumn>
          <TableColumn width={80}>{t("team.member.state.role")}</TableColumn>
          <TableColumn width={80}>Action</TableColumn>
        </TableHeader>
        <TableBody
          items={data?.data.data || []}
          loadingContent={<Spinner />}
          loadingState={isLoading ? "loading" : "idle"}
        >
          {(item) => (
            <TableRow key={`${item.id}-${item.user.id}`}>
              <TableCell>{item.id}</TableCell>
              <TableCell>{item.user.full_name}</TableCell>
              <TableCell>{item.user.email}</TableCell>
              <TableCell>{item.team_role.name}</TableCell>
              <TableCell>
                <Button isIconOnly size="sm" variant="bordered">
                  <MoreHorizontal size={18} />
                </Button>
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
