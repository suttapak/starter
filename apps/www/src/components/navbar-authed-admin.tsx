import {
  cn,
  Navbar as HeroUINavbar,
  link,
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
import { LogOut, Menu, User2 } from "lucide-react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";

import MNavbarMenuItems, { MNavbarMenuItemsProps } from "./navbar-menu-item";

import { ThemeSwitch } from "@/components/theme-switch";
import { useGetImageProfile, useGetUserMe } from "@/hooks/use-user";
import { useLogout } from "@/hooks/use-auth";
import { toastMessage } from "@/utils/toastMessage";

export type NavbarAuthedAdminProps = {
  menus: MNavbarMenuItemsProps[];
};

export const NavbarAuthedAdmin = (props: NavbarAuthedAdminProps) => {
  const { menus } = props;
  const { t } = useTranslation();
  const { isLoading: isLoadingUserMe, data: userMe } = useGetUserMe();
  const {
    onOpen: onOpenProfile,
    onClose: onCloseProfile,
    ...drawerProfileState
  } = useDisclosure();
  const {
    onOpen: onOpenMenu,
    onClose: onCloseMenu,
    ...drawerMenuState
  } = useDisclosure();
  const { mutateAsync, isPending } = useLogout();
  const { src } = useGetImageProfile();

  const profile: MNavbarMenuItemsProps[] = [
    {
      path: "/profile",
      name: "profile",
      display_name: t("navbar.profile.profile"),
      icon: <User2 size={18} />,
      padding: 0,
      onClose: onCloseProfile,
    },
  ];

  return (
    <>
      <HeroUINavbar maxWidth="full" position="sticky">
        <NavbarContent className="basis-1/5 sm:basis-full" justify="start">
          <NavbarItem>
            <Button
              isIconOnly
              size="sm"
              variant="bordered"
              onPress={onOpenMenu}
            >
              <Menu />
            </Button>
          </NavbarItem>
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
            {isLoadingUserMe ? (
              <Skeleton className="h-10 w-10 rounded-full" />
            ) : (
              <Avatar
                isIconOnly
                as={Button}
                name={userMe?.data.data.full_name}
                src={src}
                onPress={onOpenProfile}
              />
            )}
          </NavbarItem>
        </NavbarContent>
      </HeroUINavbar>
      <Drawer size="xs" {...drawerProfileState}>
        <DrawerContent>
          <DrawerHeader>
            <User
              avatarProps={{ name: userMe?.data.data.full_name, src: src }}
              description={userMe?.data.data.email}
              name={userMe?.data.data.full_name}
            />
          </DrawerHeader>
          <DrawerBody>
            <ul className="flex flex-col gap-1">
              {profile.map((item) => (
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
      <Drawer placement="left" size="xs" {...drawerMenuState}>
        <DrawerContent>
          <DrawerHeader>Admin</DrawerHeader>
          <DrawerBody>
            <ul className="flex flex-col gap-1">
              {menus.map((item) => (
                <MNavbarMenuItems
                  key={item.name}
                  {...item}
                  onClose={onCloseMenu}
                />
              ))}
            </ul>
          </DrawerBody>
          <DrawerFooter>
            <div className="flex flex-col gap-2">
              <p>{t("navbar.brand")}</p>
            </div>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
    </>
  );
};
