import { useMutation } from "@tanstack/react-query";
import axios from "axios";

import { LoginDto, RegisterDto } from "../schemas";
import { AuthResponse } from "../types/auth-response";

import { useAuth } from "@/auth";
import { Response } from "@/shared/types/api-response";
import { postJson } from "@/core/utils";

const setAxiosAuthHeader = (token: string) => {
  axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
};

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
