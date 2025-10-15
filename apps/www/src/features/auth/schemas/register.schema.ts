import { z } from "zod";

export const registerSchema = z
  .object({
    username: z.string(),
    password: z.string().min(8, ""),
    full_name: z.string().min(1, ""),
    email: z.email(),
    confirm_password: z.string().min(8, ""),
  })
  .refine((data) => data.password === data.confirm_password, {
    path: ["confirm_password"],
    error: "Passwords do not match",
  });

export type RegisterDto = z.infer<typeof registerSchema>;
