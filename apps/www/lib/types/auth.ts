import { z } from "zod";

export const LoginSchema = z.object({
  name: z
    .string()
    .min(2, { message: "Name must be at least 2 characters long." })
    .trim(),
});

export type LoginType = z.infer<typeof LoginSchema>;

export type LoginState =
  | {
      errors?: {
        name?: string[];
      };
      message?: string;
    }
  | undefined;

export const RegisterSchema = z.object({
  username: z
    .string()
    .min(2, { message: "username must be at least 2 characters long." })
    .trim(),
  email: z.string().email().trim(),
  password: z
    .string()
    .min(8, { message: "Password must be at least 8 characters long." })
    .max(32, { message: "Password must be at lessthan 32 characters long." })
    .trim()
    .refine((password) => /[A-Z]/.test(password), {
      message: "Password must have Uppercase character",
    })
    .refine((password) => /[a-z]/.test(password), {
      message: "Password must have Lowercase character",
    })
    .refine((password) => /[0-9]/.test(password), {
      message: "Password must have number",
    })
    .refine((password) => /[!@#$%^&*]/.test(password), {
      message: "Passowrd must have symbol",
    }),
  cf_password: z
    .string()
    .min(8, { message: "Password must be at least 8 characters long." })
    .max(32, { message: "Password must be at lessthan 32 characters long." })
    .trim()
    .refine((password) => /[A-Z]/.test(password), {
      message: "Password must have Uppercase character",
    })
    .refine((password) => /[a-z]/.test(password), {
      message: "Password must have Lowercase character",
    })
    .refine((password) => /[0-9]/.test(password), {
      message: "Password must have number",
    })
    .refine((password) => /[!@#$%^&*]/.test(password), {
      message: "Passowrd must have symbol",
    }),
  first_name: z
    .string()
    .min(2, { message: "first name must be at least 2 characters long." })
    .trim(),
  last_name: z
    .string()
    .min(2, { message: "first name must be at least 2 characters long." })
    .trim(),
});

export type RegisterType = z.infer<typeof RegisterSchema>;

export type RegisterState =
  | {
      errors?: {
        username?: string[];
        email?: string[];
        password?: string[];
        cf_password?: string[];
        first_name?: string[];
        last_name?: string[];
      };
      message?: string;
      data?: RegisterType;
    }
  | undefined;
