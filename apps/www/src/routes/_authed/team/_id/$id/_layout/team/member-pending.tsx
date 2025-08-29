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
  Chip,
} from "@heroui/react";
import { z } from "zod";
import { t } from "i18next";

import { useGetTeamMemberPending } from "@/hooks/use-team";
import AcceptTeamMember from "@/components/accept-team-member";

const validateSearch = z.object({
  page: z.number().prefault(1),
  limit: z.number().prefault(10),
  code: z.string().optional(),
  name: z.string().optional(),
  uom: z.string().optional(),
});

export const Route = createFileRoute(
  "/_authed/team/_id/$id/_layout/team/member-pending",
)({
  component: RouteComponent,
  validateSearch,
});

function RouteComponent() {
  const { data, isLoading } = useGetTeamMemberPending();
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
                  <option value="20">20</option>
                  <option value="50">50</option>
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
          <TableColumn width={80}>{t("team.member.state.status")}</TableColumn>
          <TableColumn width={80}>Action</TableColumn>
        </TableHeader>
        <TableBody
          emptyContent={t("not_found.global.title")}
          items={data?.data.data || []}
          loadingContent={<Spinner />}
          loadingState={isLoading ? "loading" : "idle"}
        >
          {(item) => (
            <TableRow key={`${item.id}-${item.user.id}`}>
              <TableCell>{item.id}</TableCell>
              <TableCell>{item.user.full_name}</TableCell>
              <TableCell>{item.user.email}</TableCell>
              <TableCell>
                <Chip>{item.is_active ? "Active" : "Pending"}</Chip>
              </TableCell>
              <TableCell>
                <AcceptTeamMember team={item} />
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
