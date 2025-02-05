"use client";

import { create } from "zustand";
import { immer } from "zustand/middleware/immer";
import { createJSONStorage, devtools, persist } from "zustand/middleware";

import { TeamType } from "@/hooks/use-team";
import useStore from "@/hooks/use-store";

type State = {
  team?: TeamType;
};
type Action = {
  setTeam: (team: TeamType) => void;
};

type TeamStore = State & Action;

export const teamStore = create<TeamStore>()(
  devtools(
    persist(
      immer((set) => ({
        team: undefined,
        setTeam: (team) =>
          set((state) => {
            state.team = team;
          }),
      })),
      {
        name: "team-storage",
        storage: createJSONStorage(() => sessionStorage),
      },
    ),
  ),
);

export const useTeamStateStore = () => {
  return useStore(teamStore, (state) => state);
};
