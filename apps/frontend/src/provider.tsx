import {
  useRouter,
  type NavigateOptions,
  type ToOptions,
} from "@tanstack/react-router";
import { NextUIProvider } from "@nextui-org/react";

declare module "@react-types/shared" {
  interface RouterConfig {
    href: ToOptions["to"];
    routerOptions: Omit<NavigateOptions, keyof ToOptions>;
  }
}

export function Provider({ children }: { children: React.ReactNode }) {
  let router = useRouter();

  return (
    <NextUIProvider
      navigate={(to, options) => router.navigate({ to, ...options })}
      useHref={(to) => router.buildLocation({ to }).href}
    >
      {children}
    </NextUIProvider>
  );
}
