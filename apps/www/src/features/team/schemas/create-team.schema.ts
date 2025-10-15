import { z } from "zod";

export const createTeamSchema = z.object({
  name: z.string().min(1, "Name is required"),
  username: z.string().min(1, "Username is required"),
  description: z.string().optional(),
});

export type CreateTeamDto = z.infer<typeof createTeamSchema>;
