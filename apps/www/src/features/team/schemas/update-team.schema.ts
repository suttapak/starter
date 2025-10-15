import { z } from "zod";

export const updateTeamSchema = z.object({
  name: z.string().optional(),
  username: z.string().optional(),
  address: z.string().optional(),
  phone: z.string().optional(),
  email: z.string().optional(),
  description: z.string().optional(),
});

export type UpdateTeamDto = z.infer<typeof updateTeamSchema>;
