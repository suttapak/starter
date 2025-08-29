import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { z } from "zod";
import { useNavigate, useParams, useSearch } from "@tanstack/react-router";

import { PaginatedResponse, Response } from "@/types/api_response";
import { getJson, postJson, putJson } from "@/utils/fetch";
import { TeamResponse } from "@/types/team_response";
import { TeamMemberResponse } from "@/types/team_member_response";

export const keys = {
  team: ["team", "me"] as const,
  teamId: (id: string) => ["team", "id", id] as const,
  member: (id: string) => ["team", "member", id] as const,
  pending: (id: string) => ["team", "member", "pending", id] as const,
  search: (page: number, limit: number, name?: string) =>
    ["team", "search", name, page, limit] as const,
};

export const useGetTeamMe = () => {
  return useQuery({
    queryKey: keys.team,
    queryFn: () => getJson<Response<TeamResponse[]>>("/teams/me"),
  });
};

export const createTeamSchema = z.object({
  name: z.string().min(1, "Name is required"),
  username: z.string().min(1, "Username is required"),
  description: z.string().optional(),
});

export type CreateTeamDto = z.infer<typeof createTeamSchema>;

export const useCreateTeam = (onSuccess?: () => void) => {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTeamDto) =>
      postJson<Response<TeamResponse>>("/teams/", data),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: keys.team });
      onSuccess?.();
    },
  });
};

export const useGetTeamById = () => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  return useQuery({
    queryKey: keys.teamId(id),
    queryFn: () => getJson<Response<TeamResponse>>(`/teams/${id}`),
    enabled: !!id,
  });
};

export const useGetTeamMembers = () => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  return useQuery({
    queryKey: keys.member(id),
    queryFn: () =>
      getJson<PaginatedResponse<TeamMemberResponse>>(`/teams/${id}/members`),
    enabled: !!id,
  });
};

export const useGetTeamMemberPending = () => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  return useQuery({
    queryKey: keys.pending(id),
    queryFn: () =>
      getJson<PaginatedResponse<TeamMemberResponse>>(
        `/teams/${id}/pending-members`,
      ),
    enabled: !!id,
  });
};

export const useShareTeam = (onSuccess?: (link: string) => void) => {
  const client = useQueryClient();
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });

  return useMutation({
    mutationFn: () =>
      postJson<Response<string>>(`/teams/${id}/shared-link`, {}),
    onSuccess: (data) => {
      client.invalidateQueries({ queryKey: keys.member(id) });
      onSuccess?.(data.data.data);
    },
  });
};

export const useJoinTeam = (onSuccess?: () => void) => {
  const client = useQueryClient();
  const { token } = useSearch({ from: "/_authed/team/_team/join-team" });
  const navigate = useNavigate();

  return useMutation({
    mutationFn: () =>
      postJson<Response<TeamMemberResponse>>(
        `/teams/join/link`,
        {},
        { params: { token } },
      ),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: keys.team });
      onSuccess?.();
      navigate({ to: "/team" });
    },
  });
};

export const useSearchTeam = () => {
  const { name, page, limit } = useSearch({
    from: "/_authed/team/_team/search",
  });

  return useQuery({
    queryKey: keys.search(page, limit, name),
    queryFn: () =>
      getJson<PaginatedResponse<TeamResponse>>(`/teams/`, {
        page,
        limit,
        name,
      }),
  });
};

export const useRequestJoinTeam = (onSuccess?: () => void) => {
  return useMutation({
    mutationFn: (teamId: number) =>
      postJson<Response<never>>(`/teams/${teamId}/request-join`, {}),
    onSuccess: () => {
      onSuccess?.();
    },
  });
};

export const acceptTeamMemberSchema = z.object({
  user_id: z.number(),
  role_id: z.number(),
});

export type AcceptTeamMemberDto = z.infer<typeof acceptTeamMemberSchema>;

export const useAcceptTeamMember = (onSuccess?: () => void) => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });
  const client = useQueryClient();

  return useMutation({
    mutationFn: (data: AcceptTeamMemberDto) =>
      postJson(`/teams/${id}/accept`, data),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: keys.member(id) });
      client.invalidateQueries({ queryKey: keys.pending(id) });
      onSuccess?.();
    },
  });
};

export const updateTeamSchema = z.object({
  name: z.string().optional(),
  username: z.string().optional(),
  address: z.string().optional(),
  phone: z.string().optional(),
  email: z.string().optional(),
  description: z.string().optional(),
});

export type UpdateTeamDto = z.infer<typeof updateTeamSchema>;

export const useUpdateTeamInfo = (onSuccess?: () => void) => {
  const { id } = useParams({ from: "/_authed/team/_id/$id/_layout" });
  const client = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateTeamDto) => putJson(`/teams/${id}`, data),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: keys.teamId(id) });
      onSuccess?.();
    },
  });
};
