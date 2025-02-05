import Link from "next/link";
import { ShieldAlert } from "lucide-react";

import { Button } from "@/components/ui/button";

export default function Forbidden() {
  return (
    <div className="flex h-[calc(100vh-4rem)] flex-col items-center justify-center gap-4">
      <div className="flex items-center gap-2">
        <ShieldAlert className="h-10 w-10 text-muted-foreground" />
        <h2 className="text-2xl font-bold">403 Forbidden</h2>
      </div>
      <p className="text-muted-foreground">
        You don&apos;t have permission to access this page.
      </p>
      <Button asChild variant="outline">
        <Link href="/">Return Home</Link>
      </Button>
    </div>
  );
}
