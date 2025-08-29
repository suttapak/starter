import { createFileRoute, Outlet } from "@tanstack/react-router";
import { cn } from "@heroui/react";
import { useTranslation } from "react-i18next";
import { UserLock } from "lucide-react";

import { NavbarAuthedAdmin } from "@/components/navbar-authed-admin";
import MNavbarMenuItems, {
  MNavbarMenuItemsProps,
} from "@/components/navbar-menu-item";

export const Route = createFileRoute("/_authed/admin/_layout")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();

  const menus: MNavbarMenuItemsProps[] = [
    {
      path: `/admin`,
      name: "admin",
      display_name: t("navbar.admin.admin"),
      icon: <UserLock size={18} />,
      padding: 0,
    },
    {
      path: `/admin/report`,
      name: "report",
      display_name: t("navbar.admin.report"),
      icon: <UserLock size={18} />,
      padding: 0,
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
          {t("navbar.admin.header_title")}
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
        <NavbarAuthedAdmin menus={menus} />
        <main className="container mx-auto max-w-7xl flex-grow p-2">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
