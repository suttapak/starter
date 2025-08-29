import { cn, link } from "@heroui/react";
import {
  Navbar as HeroUINavbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
} from "@heroui/react";
import {
  Avatar,
  Button,
  Drawer,
  DrawerBody,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  Skeleton,
  useDisclosure,
  User,
} from "@heroui/react";
import { LogOut, User2 } from "lucide-react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";

import MNavbarMenuItems, { MNavbarMenuItemsProps } from "./navbar-menu-item";

import { ThemeSwitch } from "@/components/theme-switch";
import { useGetImageProfile, useGetUserMe } from "@/hooks/use-user";
import { useLogout } from "@/hooks/use-auth";
import { toastMessage } from "@/utils/toastMessage";

export const NavbarAuthed = () => {
  const { t } = useTranslation();
  const { isLoading, data } = useGetUserMe();
  const { onOpen, onClose, ...drawerState } = useDisclosure();
  const { mutateAsync, isPending } = useLogout();
  const { src } = useGetImageProfile();

  const navMenu: MNavbarMenuItemsProps[] = [
    {
      path: "/profile",
      name: "profile",
      display_name: t("navbar.profile.profile"),
      icon: <User2 size={18} />,
      padding: 0,
      onClose: onClose,
    },
  ];

  return (
    <>
      <HeroUINavbar maxWidth="full" position="sticky">
        <NavbarContent className="basis-1/5 sm:basis-full" justify="start">
          <NavbarBrand className="gap-3 max-w-fit">
            <Link
              className={cn(
                "flex justify-start items-center gap-1",
                link({ color: "foreground" }),
              )}
              to="/team"
            >
              <p className="font-bold text-inherit">{t("navbar.brand")}</p>
            </Link>
          </NavbarBrand>
        </NavbarContent>

        <NavbarContent justify="end">
          <ThemeSwitch />
          <NavbarItem>
            {isLoading ? (
              <Skeleton className="h-10 w-10 rounded-full" />
            ) : (
              <Avatar
                isIconOnly
                as={Button}
                name={data?.data.data.full_name}
                src={src}
                onPress={onOpen}
              />
            )}
          </NavbarItem>
        </NavbarContent>
      </HeroUINavbar>
      <Drawer size="xs" {...drawerState}>
        <DrawerContent>
          <DrawerHeader>
            <User
              avatarProps={{ name: data?.data.data.full_name, src: src }}
              description={data?.data.data.email}
              name={data?.data.data.full_name}
            />
          </DrawerHeader>
          <DrawerBody>
            <ul className="flex flex-col gap-1">
              {navMenu.map((item) => (
                <MNavbarMenuItems key={item.name} {...item} />
              ))}
            </ul>
          </DrawerBody>
          <DrawerFooter>
            <Button
              fullWidth
              color="danger"
              endContent={<LogOut size={18} />}
              isLoading={isPending}
              size="sm"
              variant="bordered"
              onPress={() => toast.promise(() => mutateAsync(), toastMessage)}
            >
              {t("navbar.profile.logout")}
            </Button>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
    </>
  );
};
