import { cn, link } from "@heroui/react";
import {
  Navbar as HeroUINavbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
} from "@heroui/react";
import { Button } from "@heroui/react";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";

import { ThemeSwitch } from "./theme-switch";

export const NavbarPublic = () => {
  const { t } = useTranslation();

  return (
    <HeroUINavbar maxWidth="full" position="sticky">
      <NavbarContent className="basis-1/5 sm:basis-full" justify="start">
        <NavbarBrand className="gap-3 max-w-fit">
          <Link
            className={cn(
              "flex justify-start items-center gap-1",
              link({ color: "foreground" }),
            )}
            to="/"
          >
            <p className="font-bold text-inherit">{t("navbar.brand")}</p>
          </Link>
        </NavbarBrand>
      </NavbarContent>

      <NavbarContent className=" pl-4" justify="end">
        <ThemeSwitch />
        <NavbarItem>
          <Button as={Link} color="primary" to="/profile">
            {t("navbar.use_app")}
          </Button>
        </NavbarItem>
      </NavbarContent>
    </HeroUINavbar>
  );
};
