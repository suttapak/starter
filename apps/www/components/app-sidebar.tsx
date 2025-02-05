"use client";

import * as React from "react";
import { Bot, Settings2, SquareTerminal } from "lucide-react";

import { NavMain } from "@/components/nav-main";
import { NavUser } from "@/components/nav-user";
import { TeamSwitcher } from "@/components/team-switcher";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";

// This is sample data.
const data = {
  navMain: [
    {
      title: "การจัดการสินค้า",
      url: "/stock-service/products",
      icon: SquareTerminal,
      isActive: true,
      items: [
        {
          title: "รายการสินค้า",
          url: "/stock-service/products",
        },
        {
          title: "เพิ่มสินค้า",
          url: "/stock-service/products/new",
        },
        {
          title: "คลังสินค้า",
          url: "/stock-service/warehouse",
        },
      ],
    },
    {
      title: "การจัดการรายการซื้อขาย",
      url: "#",
      icon: Bot,
      isActive: true,
      items: [
        {
          title: "รายการซื้อขาย",
          url: "#",
        },
      ],
    },

    {
      title: "การตั้งค้าคลัง",
      url: "#",
      icon: Settings2,
      isActive: true,
      items: [
        {
          title: "ทัวไป",
          url: "#",
        },
        {
          title: "บุครากร",
          url: "#",
        },
        {
          title: "การเชื่อมต่อภายนอก",
          url: "#",
        },
      ],
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <TeamSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
