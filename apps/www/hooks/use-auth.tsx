import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

import { postJson } from "@/lib/api";

export const useLogout = () => {
  const router = useRouter();

  return useMutation({
    mutationFn: () => postJson("/auth/logout"),
    onSuccess: () => {
      router.replace("/login");
    },
  });
};
