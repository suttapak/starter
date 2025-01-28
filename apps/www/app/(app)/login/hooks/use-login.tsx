"use client";
import { useMutation } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";

import { LoginDto } from "../ui/login-form";

import { postJson } from "@/lib/api";

export const useLogin = () => {
  const router = useRouter();
  const searchParams = useSearchParams(); // Access query parameters

  return useMutation({
    mutationFn: (data: LoginDto) => {
      return postJson("/auth/login", data);
    },
    onSuccess: () => {
      const redirectPath = searchParams.get("redirect") || "/stock-service"; // Default to "/stock-service"

      router.replace(redirectPath); // Navigate to the redirect path
    },
  });
};
