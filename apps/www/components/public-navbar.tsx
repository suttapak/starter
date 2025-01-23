"use client";
import React from "react";
import {
  Navbar as HeroNavbar,
  NavbarBrand,
  NavbarContent,
  Divider,
  Button,
} from "@heroui/react";
import Link from "next/link";
import { LogIn, UserPlus } from "lucide-react";

import { ThemeSwitch } from "./theme-switch";

function PublicNavbar() {
  return (
    <>
      <HeroNavbar
        classNames={{
          wrapper: "max-w-8xl",
        }}
        maxWidth="xl"
        position="sticky"
      >
        <NavbarContent className="sm:hidden">
          <NavbarBrand>
            <p className="font-bold text-inherit">Labotron</p>
          </NavbarBrand>
        </NavbarContent>

        <NavbarContent justify="end">
          <ThemeSwitch />
          <Divider className="h-7 hidden sm:flex" orientation="vertical" />
          <Button as={Link} endContent={<UserPlus />} href="/register">
            Register
          </Button>
          <Button
            as={Link}
            color="primary"
            endContent={<LogIn />}
            href="/login"
          >
            Login
          </Button>
        </NavbarContent>
      </HeroNavbar>
    </>
  );
}

export default PublicNavbar;
