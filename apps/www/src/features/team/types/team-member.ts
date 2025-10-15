import type { CommonModel } from "@/shared/types";
import type { UserResponse } from "@/features/user";

export interface TeamMemberResponse extends CommonModel {
  team_id: number;
  user_id: number;
  team_role_id: number;
  user: UserResponse;
  team_role: TeamRoleResponse;
  is_active: boolean;
}

export interface TeamRoleResponse extends CommonModel {
  name: string;
}
