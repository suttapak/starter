import { useMutation } from "@tanstack/react-query";
import { z } from "zod";
import axios from "axios";

import { postJson } from "@/utils/fetch";
import { Response } from "@/types/api_response";
import { AuthResponse } from "@/types/auth_response";
import { useAuth } from "@/auth";

export const loginSchema = z.object({
  username: z.string(),
  password: z.string().min(8, ""),
});

const setAxiosAuthHeader = (token: string) => {
  axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
};

export type LoginDto = z.infer<typeof loginSchema>;

export const useLogin = (onSuccess?: () => Promise<void>) => {
  return useMutation({
    mutationFn: async (body: LoginDto) => {
      const res = await postJson<Response<AuthResponse>>("/auth/login", body);

      return res;
    },
    onSuccess: async ({ data }) => {
      localStorage.setItem("accessToken", data.data.token);
      setAxiosAuthHeader(data.data.token);
      localStorage.setItem("refreshToken", data.data.refresh_token);
      await onSuccess?.();
    },
  });
};

export const useRefreshToken = (onSuccess?: () => void) => {
  return useMutation({
    mutationFn: () => {
      const refreshToken = localStorage.getItem("refreshToken");

      if (!refreshToken) {
        throw new Error("No refresh token found");
      }
      setAxiosAuthHeader(refreshToken);

      return postJson<Response<AuthResponse>>("/auth/refresh", {});
    },
    onSuccess: async ({ data }) => {
      localStorage.setItem("accessToken", data.data.token);
      setAxiosAuthHeader(data.data.token);
      localStorage.setItem("refreshToken", data.data.refresh_token);
      onSuccess?.();
    },
  });
};

export const useLogout = (onSuccess?: () => void) => {
  const { onChangeIsAuthenticated } = useAuth();

  return useMutation({
    mutationFn: () => postJson<never>("/auth/logout", {}),
    onSuccess: () => {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      onChangeIsAuthenticated(false);
      onSuccess?.();
    },
  });
};

export const useSendVerifyEmail = (onSuccess?: () => void) => {
  return useMutation({
    mutationFn: () => postJson<never>("/auth/email/send-verify", {}),
    onSuccess: () => {
      onSuccess?.();
    },
  });
};
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

export const useRegister = (onSuccess?: () => void) => {
  return useMutation({
    mutationFn: (body: RegisterDto) =>
      postJson<Response<AuthResponse>>("/auth/register", body),
    onSuccess: async ({ data }) => {
      localStorage.setItem("accessToken", data.data.token);
      localStorage.setItem("refreshToken", data.data.refresh_token);
      onSuccess?.();
    },
  });
};
