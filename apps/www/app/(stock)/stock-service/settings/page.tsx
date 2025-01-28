"use client";
import { Moon } from "lucide-react";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ModeToggle } from "@/components/mode-toggle";

function Page() {
  return (
    <div className="flex justify-center">
      <Card className="w-full max-w-screen-sm">
        <CardHeader>
          <CardTitle>สิ่งที่ปรากฎ</CardTitle>
          <CardDescription>ปรับแต่งลักษณะของแอปบนอุปกรณ์ของคุณ</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Moon className="h-5 w-5" />

              <div className="space-y-0.5">
                <p className="font-medium">Dark Mode</p>
                <p className="text-sm text-muted-foreground">
                  สลับระหว่างธีมสว่างและธีมมืด
                </p>
              </div>
            </div>
            <ModeToggle />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export default Page;
