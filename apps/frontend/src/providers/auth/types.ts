export type AuthFunction = {
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

export type AuthState = {
  user: TypeGetusersByID | null;
  isInit: boolean;
  isAuthenticated: boolean;
};

export type TypeGetusersByID = {};
