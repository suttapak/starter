import React, { PropsWithChildren } from "react";

import { SettingsTabs } from "@/components/settings-tabs";
import { title } from "@/components/primitives";

type Props = PropsWithChildren & {};

function Layout({ children }: Props) {
  return (
    <div className="container">
      <div className="py-8">
        <h1 className={title({ size: "sm" })}>ตั่งค่าข้อมูลผู้ใช้งาน</h1>
      </div>
      <SettingsTabs />
      <section className="py-2">{children}</section>
    </div>
  );
}

export default Layout;
