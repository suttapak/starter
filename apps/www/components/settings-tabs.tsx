// SettingsTabs.tsx
"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";

import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";

export const SettingsTabs = () => {
  const path = usePathname();

  return (
    <Tabs className="w-full" value={path}>
      <TabsList className="w-full p-0 bg-background justify-start border-b rounded-none">
        <TabsLinkTrigger href="/stock-service/settings">
          Settings
        </TabsLinkTrigger>
        <TabsLinkTrigger href="/stock-service/settings/account">
          Account
        </TabsLinkTrigger>
        <TabsLinkTrigger href="/stock-service/settings/notification">
          Notification
        </TabsLinkTrigger>
      </TabsList>
    </Tabs>
  );
};

const TabsLinkTrigger: React.FC<{
  href: string;
  children: React.ReactNode;
}> = ({ href, children }) => (
  <TabsTrigger
    asChild
    className="rounded-none bg-background h-full data-[state=active]:shadow-none border-b-2 border-transparent data-[state=active]:border-primary"
    value={href}
  >
    <Link href={href}>{children}</Link>
  </TabsTrigger>
);
