"use client";
import { Button, Link } from "@heroui/react";
import { ChevronDown } from "lucide-react";
import React from "react";

import { ScrollArea } from "./scroll-area";

function Sidebar() {
  return (
    <ScrollArea className="h-full max-h-[calc(100vh_-_64px)]">
      <div className="flex flex-col p-2">
        <Button
          fullWidth
          as={Link}
          className="justify-between font-semibold"
          endContent={<ChevronDown />}
          href="/"
          size="sm"
          variant="light"
        >
          Home
        </Button>
      </div>
    </ScrollArea>
  );
}

export default Sidebar;
