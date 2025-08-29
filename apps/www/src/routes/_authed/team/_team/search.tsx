import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import {
  Input,
  Pagination,
  Spinner,
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
} from "@heroui/react";
import { useTranslation } from "react-i18next";

import { useSearchTeam } from "@/hooks/use-team";
import TeamRequestJoinModal from "@/components/team-request-join-modal";
const validateSearch = z.object({
  name: z.string().optional(),
  page: z.number().prefault(1),
  limit: z.number().prefault(10),
});

export const Route = createFileRoute("/_authed/team/_team/search")({
  component: RouteComponent,
  validateSearch,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { data, isLoading } = useSearchTeam();
  const navigate = Route.useNavigate();
  const { page, limit, name } = Route.useSearch();

  return (
    <div className="flex flex-col gap-2">
      <Input
        placeholder={t("team.search")}
        value={name}
        onValueChange={(value) => navigate({ search: { name: value } })}
      />
      <Table
        aria-label="ตารางแสดงข้อมูลแผนก"
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
          <TableColumn width={80}>{t("team.state.index")}</TableColumn>
          <TableColumn width={200}>{t("team.state.name")}</TableColumn>
          <TableColumn width={200}>{t("team.state.username")}</TableColumn>
          <TableColumn width={80}>{t("team.state.description")}</TableColumn>
          <TableColumn width={80}>Action</TableColumn>
        </TableHeader>
        <TableBody
          emptyContent={name ? `ไม่พบแผนกที่ค้นหา [${name}]` : "ไม่มีแผนก"}
          items={data?.data.data || []}
          loadingContent={<Spinner />}
          loadingState={isLoading ? "loading" : "idle"}
        >
          {(item) => (
            <TableRow key={`${item.id}`}>
              <TableCell>{item.id}</TableCell>
              <TableCell>{item.username}</TableCell>
              <TableCell>{item.name}</TableCell>
              <TableCell>{item.description}</TableCell>
              <TableCell>
                <TeamRequestJoinModal team={item} />
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
