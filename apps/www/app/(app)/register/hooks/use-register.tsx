import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

import { RegisterDto } from "../ui/register-form";

import { postJson } from "@/lib/api";

export const useRegister = () => {
  const router = useRouter();

  return useMutation({
    mutationFn: (data: RegisterDto) => {
      return postJson("/auth/register", data);
    },
    onSuccess: () => {
      router.replace("/login");
    },
  });
};
