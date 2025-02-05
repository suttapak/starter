import Link from "next/link";
import { Lock } from "lucide-react";

import { Button } from "@/components/ui/button";

export default function Unauthorized() {
  return (
    <div className="flex h-[calc(100vh-4rem)] flex-col items-center justify-center gap-4">
      <div className="flex items-center gap-2">
        <Lock className="h-10 w-10 text-muted-foreground" />
        <h2 className="text-2xl font-bold">401 Unauthorized</h2>
      </div>
      <p className="text-muted-foreground">
        Please log in to access this page.
      </p>
      <div className="flex gap-4">
        <Button asChild variant="outline">
          <Link href="/">Return Home</Link>
        </Button>
        <Button asChild>
          <Link href="/login">Log In</Link>
        </Button>
      </div>
    </div>
  );
}
