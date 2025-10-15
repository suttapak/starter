import type { PaginatedResponse, Response } from "@/shared/types";
import type { TeamResponse } from "../types/team";
import type { TeamMemberResponse } from "../types/team-member";
import type { CreateTeamDto } from "../schemas/create-team.schema";
import type { UpdateTeamDto } from "../schemas/update-team.schema";
import type { AcceptTeamMemberDto } from "../schemas/accept-team-member.schema";

import { useNavigate, useParams, useSearch } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { getJson, postJson, putJson } from "@/core/utils/fetch";

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
