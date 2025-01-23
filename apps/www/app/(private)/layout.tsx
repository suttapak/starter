import "@/styles/globals.css";
import { Metadata, Viewport } from "next";
import clsx from "clsx";
import { Divider } from "@heroui/react";

import { siteConfig } from "@/config/site";
import { fontSans } from "@/config/fonts";
import Navbar from "@/components/navbar";
import Sidebar from "@/components/sidebar";

export const metadata: Metadata = {
  title: {
    default: siteConfig.name,
    template: `%s - ${siteConfig.name}`,
  },
  description: siteConfig.description,
  icons: {
    icon: "/favicon.ico",
  },
};

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "white" },
    { media: "(prefers-color-scheme: dark)", color: "black" },
  ],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html suppressHydrationWarning lang="en">
      <head />
      <body
        className={clsx(
          "min-h-screen bg-background font-sans antialiased",
          fontSans.variable,
        )}
      >
        <div className="hidden sm:block bg-content2 w-60 fixed h-dvh z-50 rounded-large rounded-l-none shadow">
          <header className="text-large font-semibold py-4 px-6">
            <h1 className="pt-[3px]">Labotron</h1>
          </header>
          <Divider />
          <Sidebar />
        </div>
        <div className="sm:ml-60">
          <Navbar />
          <main className="px-6 py-4">{children}</main>
        </div>
      </body>
    </html>
  );
}
