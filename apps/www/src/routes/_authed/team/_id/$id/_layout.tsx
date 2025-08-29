import { createFileRoute, Outlet, useParams } from "@tanstack/react-router";
import { cn, Skeleton, User } from "@heroui/react";
import {
  ChartCandlestick,
  Warehouse,
  Plus,
  ChartColumnStacked,
  Contact,
  Users,
  UserRoundPlus,
} from "lucide-react";
import { useTranslation } from "react-i18next";

import { NavbarAuthedTeam } from "@/components/navbar-authed-team";
import MNavbarMenuItems, {
  MNavbarMenuItemsProps,
} from "@/components/navbar-menu-item";
import { useGetTeamById } from "@/hooks/use-team";

export const Route = createFileRoute("/_authed/team/_id/$id/_layout")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { isLoading: isLoadingTeam, data: team } = useGetTeamById();
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  const menus: MNavbarMenuItemsProps[] = [
    {
      path: `/team/${id}`,
      name: "summary",
      display_name: t("navbar.team.report"),
      icon: <ChartCandlestick size={18} />,
      padding: 0,
    },
    {
      name: "product",
      display_name: t("navbar.team.stock.title"),
      icon: <Warehouse size={18} />,
      padding: 0,
      child: [
        {
          path: `/team/${id}/product`,
          name: "product-list",
          display_name: t("navbar.team.stock.stock"),
          icon: <Warehouse size={18} />,
          padding: 0,
        },
        {
          path: `/team/${id}/product/new`,
          name: "product-new",
          display_name: t("navbar.team.stock.new_product"),
          icon: <Plus size={18} />,
          padding: 0,
        },
        {
          path: `/team/${id}/product/category`,
          name: "product-category",
          display_name: t("navbar.team.stock.category"),
          icon: <ChartColumnStacked size={18} />,
          padding: 0,
        },
      ],
    },
    {
      name: "team",
      display_name: t("navbar.team.info.title"),
      icon: <Contact size={18} />,
      padding: 0,
      child: [
        {
          path: `/team/${id}/team`,
          name: "team/info",
          display_name: t("navbar.team.info.info"),
          icon: <Contact size={18} />,
          padding: 0,
        },
        {
          path: `/team/${id}/team/member`,
          name: "team/member",
          display_name: t("navbar.team.info.member"),
          icon: <Users size={18} />,
          padding: 0,
        },
        {
          path: `/team/${id}/team/member-pending`,
          name: "team/member-pending",
          display_name: t("navbar.team.info.member_pending"),
          icon: <UserRoundPlus size={18} />,
          padding: 0,
        },
      ],
    },
  ];

  return (
    <div className="flex h-dvh w-full overflow-hidden">
      <div
        className={cn([
          "bg-default-100 dark:bg-background",
          "border-r-1 border-transparent dark:border-default-100 ",
          "ease-sidebar-collapse overflow-hidden p-4 transition-[width] duration-200",
          "sm:flex flex-col relative hidden h-full ",
          "max-w-[16rem] w-full",
        ])}
      >
        <header className="flex py-4 flex-initial text-large font-semibold">
          {isLoadingTeam && <Skeleton className="h-10 rounded" />}
          <User
            avatarProps={{ name: team?.data.data.name }}
            description={team?.data.data.username}
            name={team?.data.data.name}
          />
        </header>
        <div className="flex flex-1 flex-col gap-3 py-2 overflow-y-auto">
          <ul className="flex flex-col gap-1">
            {menus.map((item) => (
              <MNavbarMenuItems key={item.name} {...item} />
            ))}
          </ul>
        </div>
        <footer className="flex flex-row gap-2 py-4 justify-end">
          <div className="flex flex-col gap-2">
            <p>{t("navbar.brand")}</p>
          </div>
        </footer>
      </div>
      <div className="flex h-dvh w-full flex-col overflow-auto">
        <NavbarAuthedTeam menus={menus} />
        <main className="container mx-auto max-w-7xl flex-grow p-2">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
