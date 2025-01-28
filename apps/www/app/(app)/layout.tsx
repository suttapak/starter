import Link from "next/link";
import { LogIn, UserPlus } from "lucide-react";

import { Button } from "@/components/ui/button";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { ModeToggle } from "@/components/mode-toggle";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider>
      <SidebarInset>
        <header className="flex sticky top-0 bg-background h-16 shrink-0 items-center gap-2 border-b px-4 justify-between">
          <div className="flex" />
          <div className="flex gap-2">
            <ModeToggle />
            <Button asChild variant={"secondary"}>
              <Link href={"/register"}>
                Register
                <UserPlus />
              </Link>
            </Button>
            <Button asChild variant={"default"}>
              <Link href={"/login"}>
                Login
                <LogIn />
              </Link>
            </Button>
          </div>
        </header>
        <main>{children}</main>
      </SidebarInset>
    </SidebarProvider>
  );
}
