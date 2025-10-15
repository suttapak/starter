import { z } from "zod";

export const acceptTeamMemberSchema = z.object({
  user_id: z.number(),
  role_id: z.number(),
});

export type AcceptTeamMemberDto = z.infer<typeof acceptTeamMemberSchema>;
