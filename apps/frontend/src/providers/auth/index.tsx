import { createContext, useState } from "react";
import axios from "axios";

import { AuthFunction, AuthState, TypeGetusersByID } from "./types";

import { loginService } from "@/services/auth";

export type AuthContextType = AuthFunction &
  AuthState & {
    profile: () => Promise<void>;
  };

const setSession = (accessToken: string) => {
  if (accessToken) {
    window.localStorage.setItem("accessToken", accessToken);
    axios.defaults.headers.common["Authorization"] = `Bearer ${accessToken}`;
  } else {
    window.localStorage.removeItem("accessToken");
    delete axios.defaults.headers.common["Authorization"];
  }
};

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);
export type AuthProviderProps = {
  children: React.ReactNode;
};
export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<TypeGetusersByID | null>(null);
  const [isInit, setIsInit] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const login = async (username: string, password: string) => {
    const response = await loginService(username, password);

    if (response.data) {
      const { data } = response.data;
      const { token, refresh_token } = data;

      try {
        setSession(token);
        localStorage.setItem("refreshToken", refresh_token);

        setIsAuthenticated(true);
      } catch (_) {
        throw new Error("Invalid username or password");
      }
    } else {
      throw new Error("Invalid username or password");
    }
  };
};
