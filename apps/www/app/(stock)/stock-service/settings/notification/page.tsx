"use client";
import { AlertCircle } from "lucide-react";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

function Page() {
  return (
    <div className="flex justify-center">
      <Card className="w-full max-w-screen-sm">
        <CardHeader className="text-center pb-2">
          <div className="flex justify-center mb-4">
            <AlertCircle className="h-12 w-12 text-muted-foreground" />
          </div>
          <CardTitle className="text-2xl sm:text-3xl">
            เพจไม่พร้อมใช้งาน
          </CardTitle>
        </CardHeader>
        <CardContent className="text-center pb-2">
          <p className="text-muted-foreground">
            หน้าที่คุณกำลังค้นหาอาจถูกลบไปแล้วและมีชื่ออยู่ มีการเปลี่ยนแปลง
            หรือไม่สามารถใช้งานได้ชั่วคราว
          </p>
        </CardContent>
        <CardFooter className="flex justify-center">
          <Button asChild>
            <Link href="/stock-service/">Return Home</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}

export default Page;
