"use client";

import { useEffect } from "react";
import { AlertCircle } from "lucide-react";

import { Button } from "@/components/ui/button";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    // Log the error to an error reporting service
    // eslint-disable-next-line no-console
    console.error(error);
  }, [error]);

  return (
    <div className="flex h-[calc(100vh-4rem)] flex-col items-center justify-center gap-4">
      <div className="flex items-center gap-2 text-destructive">
        <AlertCircle className="h-6 w-6" />
        <h2 className="text-lg font-semibold">Something went wrong!</h2>
      </div>
      <p className="text-muted-foreground">
        An unexpected error occurred. Please try again later.
      </p>
      <Button variant="outline" onClick={reset}>
        Try again
      </Button>
      {process.env.NODE_ENV === "development" && (
        <div className="max-w-lg rounded-md bg-muted/50 p-4">
          <p className="text-sm text-muted-foreground">{error.message}</p>
        </div>
      )}
    </div>
  );
}
