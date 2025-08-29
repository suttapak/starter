import type { NavigateOptions, ToOptions } from "@tanstack/react-router";

import { Toaster } from "react-hot-toast";
import { HeroUIProvider } from "@heroui/react";
import { PropsWithChildren } from "react";
import { useRouter } from "@tanstack/react-router";

import "@/index.css";

type Props = PropsWithChildren & {};

declare module "@react-types/shared" {
  interface RouterConfig {
    href: ToOptions["to"];
    routerOptions: Omit<NavigateOptions, keyof ToOptions>;
  }
}

const UiProvider = ({ children }: Props) => {
  let router = useRouter();

  return (
    <HeroUIProvider
      locale="en-GB"
      navigate={(to, options) => router.navigate({ to, ...options })}
      useHref={(to) => router.buildLocation({ to }).href}
    >
      {children}
      <Toaster />
    </HeroUIProvider>
  );
};

export default UiProvider;
