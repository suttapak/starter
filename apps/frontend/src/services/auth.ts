import { postJson } from "@/utils";

export type Login = {
  token: string;
  refresh_token: string;
};

export const LOGIN_URL = "/auth/login";
export const loginService = async (username: string, password: string) => {
  return await postJson<Login>(LOGIN_URL, { username, password });
};
