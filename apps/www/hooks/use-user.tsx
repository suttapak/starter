import { useQuery } from "@tanstack/react-query";

import { getJson } from "@/lib/api";
import { Response } from "@/types/response";

export interface UserRes {
  id: number;
  created_at: Date;
  updated_at: Date;
  username: string;
  email: string;
  email_verifyed: boolean;
  full_name: string;
  role_id: number;
}

const keys = {
  me: ["me"] as const,
};

export const useGetUserMe = () => {
  return useQuery({
    queryKey: keys.me,
    queryFn: () => getJson<Response<UserRes>>("/users/me"),
  });
};
