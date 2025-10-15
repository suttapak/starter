import { useState } from "react";
import {
  ErrorComponentProps,
  useNavigate,
  useRouter,
} from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { AlertTriangle, ArrowLeft, Home } from "lucide-react";
import { Button } from "@heroui/react";

const ErrorBoundary = ({ error }: ErrorComponentProps) => {
  const [message] = useState(
    isAxiosError(error) && error.response
      ? error.response.data.message || error.message
      : error.message,
  );
  const router = useRouter();
  const navigate = useNavigate();
  const handleGoBack = () => router.history.back();
  const handleGoHome = () => navigate({ to: "/" });

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background p-4">
      <div className="w-full max-w-md space-y-6 rounded-large bg-content1 p-8 shadow-sm">
        <div className="flex flex-col items-center text-center">
          <div className="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-danger/10">
            <AlertTriangle className="h-10 w-10 text-danger" />
          </div>
          <h1 className="mb-2 text-3xl font-bold text-foreground">
            Application Error
          </h1>
          <p className="mb-6 text-foreground-500">{message}</p>
        </div>

        <div className="flex flex-col gap-2 sm:flex-row">
          <Button
            fullWidth
            color="default"
            startContent={<ArrowLeft />}
            variant="flat"
            onPress={handleGoBack}
          >
            Go Back
          </Button>
          <Button
            fullWidth
            color="primary"
            startContent={<Home />}
            onPress={handleGoHome}
          >
            Go Home
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ErrorBoundary;
