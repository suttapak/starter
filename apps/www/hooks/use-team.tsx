"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { getJson, postJson } from "@/lib/api";
import { CreateTeamDto } from "@/components/new-team-dialog";
import { Response } from "@/types/response";

const keys = {
  me: ["team", "me"] as const,
};

export interface TeamType {
  id: number;
  created_at: Date;
  updated_at: Date;
  name: string;
  username: string;
  description?: string;
}

export const useGetTeamMe = () => {
  return useQuery({
    queryKey: keys.me,
    queryFn: () => getJson<Response<Array<TeamType>>>("/teams/me"),
  });
};

export const useCreateTeam = (onSuccess?: () => void) => {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTeamDto) => {
      return postJson("/teams/", data);
    },
    onSuccess: () => {
      client.invalidateQueries({ queryKey: keys.me });
      onSuccess?.();
    },
  });
};
