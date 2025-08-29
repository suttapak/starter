import { CommonModel } from "./common";
import { UserResponse } from "./user_response";

export interface TeamMemberResponse extends CommonModel {
  team_id: number;
  user_id: number;
  team_role_id: number;
  user: UserResponse;
  team_role: TeamRoleResponse;
  is_active: boolean;
}

// TeamRoleResponse interface
export interface TeamRoleResponse extends CommonModel {
  name: string;
}
