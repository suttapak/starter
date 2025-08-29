import { jwtDecode } from "jwt-decode";
import React, {
  createContext,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";

import { useGetMutateUserMe } from "./hooks/use-user";
import { useRefreshToken } from "./hooks/use-auth";

const isValidToken = (accessToken: string) => {
  if (!accessToken) {
    return false;
  }
  const decodedToken = jwtDecode(accessToken);
  const currentTime = Date.now() / 1000;

  if (!decodedToken.exp) {
    return false;
  }

  return decodedToken.exp > currentTime;
};

export type AuthContextType = {
  isInit: boolean;
  isAuthenticated: boolean;
  onChangeIsAuthenticated: (auth: boolean) => void;
};

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);

export type AuthProviderProps = {
  children: React.ReactNode;
};

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isInit, setIsInit] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const {
    mutateAsync: userMutateAsync,
    isError: userIsError,
    error: userError,
  } = useGetMutateUserMe();
  const {
    mutateAsync: refreshMutateAsync,
    isError: refreshIsError,
    error: refreshError,
  } = useRefreshToken();

  useEffect(() => {
    (async () => {
      try {
        const accessToken = localStorage.getItem("accessToken");

        if (!!accessToken && isValidToken(accessToken)) {
          await userMutateAsync();
          if (userIsError) throw userError;
          setIsAuthenticated(true);
          setIsInit(true);

          return;
        } else {
          // refresh token
          await refreshMutateAsync();
          if (refreshIsError) throw refreshError;
          setIsAuthenticated(true);
          setIsInit(true);

          return;
        }
      } catch (error) {
        setIsAuthenticated(false);
        setIsInit(true);
        throw error;
      }
    })();
  }, []);

  const onChangeIsAuthenticated = useCallback((auth: boolean) => {
    setIsAuthenticated(auth);
  }, []);

  const value = useMemo(
    () => ({
      isAuthenticated,
      isInit,
      onChangeIsAuthenticated,
    }),
    [isAuthenticated, isInit, onChangeIsAuthenticated],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = React.useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  return context;
};
