"use client";

import { MoreHorizontal, Users2 } from "lucide-react";

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface TeamCardProps {
  name: string;
  description?: string;
  username: string;
  membersCount: number;
  onEdit?: () => void;
  onDelete?: () => void;
  onLeave?: () => void;
}

export function TeamCard({
  name,
  description,
  username,
  membersCount,
  onEdit,
  onDelete,
  onLeave,
}: TeamCardProps) {
  return (
    <Card className="relative overflow-hidden transition-all hover:shadow-lg">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div className="flex items-center space-x-2">
          <Users2 className="h-5 w-5 text-muted-foreground" />
          <h3 className="font-semibold tracking-tight">{name}</h3>
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button className="h-8 w-8 p-0" variant="ghost">
              <MoreHorizontal className="h-4 w-4" />
              <span className="sr-only">Open menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={onEdit}>Edit team</DropdownMenuItem>
            <DropdownMenuItem onClick={onLeave}>Leave team</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={onDelete}
            >
              Delete team
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground line-clamp-2">
          {description}
        </p>
      </CardContent>
      <CardFooter className="flex justify-between">
        <p className="text-sm text-muted-foreground">@{username}</p>
        <div className="flex items-center space-x-1">
          <Users2 className="h-4 w-4 text-muted-foreground" />
          <span className="text-sm text-muted-foreground">{membersCount}</span>
        </div>
      </CardFooter>
    </Card>
  );
}
