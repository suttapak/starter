"use client";

import type { ThemeProviderProps } from "next-themes";

import * as React from "react";
import { ThemeProvider as NextThemesProvider } from "next-themes";
import { Toaster } from "sonner";
import { AppProgressBar as ProgressBar } from "next-nprogress-bar";

import QeuryProviders from "./query-providers";

export interface ProvidersProps {
  children: React.ReactNode;
  themeProps?: ThemeProviderProps;
}

export function Providers({ children, themeProps }: ProvidersProps) {
  return (
    <React.Fragment>
      <NextThemesProvider {...themeProps}>
        <QeuryProviders>{children}</QeuryProviders>
      </NextThemesProvider>
      <Toaster richColors position="top-center" />
      <ProgressBar color="#171717e6" options={{ showSpinner: false }} />
    </React.Fragment>
  );
}
