"use client";
import React from "react";
import {
  Navbar as HeroNavbar,
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerBody,
  DrawerFooter,
  NavbarBrand,
  NavbarContent,
  Avatar,
  Button,
  useDisclosure,
  Divider,
  Input,
} from "@heroui/react";
import { Menu, Search, SidebarClose } from "lucide-react";

import Sidebar from "./sidebar";
import { ThemeSwitch } from "./theme-switch";

function Navbar() {
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  const [isMenuOpen, setIsMenuOpen] = React.useState(false);

  return (
    <>
      <HeroNavbar
        classNames={{
          wrapper: "max-w-8xl",
        }}
        isMenuOpen={isMenuOpen}
        maxWidth="xl"
        position="sticky"
        onMenuOpenChange={setIsMenuOpen}
      >
        <NavbarContent className="sm:hidden">
          <Button isIconOnly variant="bordered" onPress={onOpen}>
            <Menu />
          </Button>
          <NavbarBrand>
            <p className="font-bold text-inherit">Labotron</p>
          </NavbarBrand>
        </NavbarContent>

        <NavbarContent justify="end">
          <Input
            classNames={{
              base: "max-w-full sm:max-w-[10rem] h-10",
              mainWrapper: "h-full",
              input: "text-small",
              inputWrapper:
                "h-full font-normal text-default-500 bg-default-400/20 dark:bg-default-500/20",
            }}
            placeholder="Type to search..."
            size="sm"
            startContent={<Search size={18} />}
            type="search"
          />
          <ThemeSwitch />
          <Divider className="h-7 hidden sm:flex" orientation="vertical" />
          <Avatar isBordered size="sm" />
        </NavbarContent>
      </HeroNavbar>
      <Drawer
        hideCloseButton
        isOpen={isOpen}
        placement={"left"}
        onOpenChange={onOpenChange}
      >
        <DrawerContent>
          {(onClose) => (
            <>
              <DrawerHeader className="flex gap-1 justify-end">
                <Button
                  isIconOnly
                  size="sm"
                  variant="bordered"
                  onPress={onClose}
                >
                  <SidebarClose />
                </Button>
              </DrawerHeader>
              <DrawerBody>
                <Sidebar />
              </DrawerBody>
              <DrawerFooter>
                <Button color="danger" variant="light" onPress={onClose}>
                  Close
                </Button>
                <Button color="primary" onPress={onClose}>
                  Action
                </Button>
              </DrawerFooter>
            </>
          )}
        </DrawerContent>
      </Drawer>
    </>
  );
}

export default Navbar;
