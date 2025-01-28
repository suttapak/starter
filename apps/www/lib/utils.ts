import { AxiosError } from "axios";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const toastMessage = {
  loading: "Loading...",
  success: "success",
  error: (error: Error) => {
    if (error instanceof AxiosError) {
      return error.response?.data?.message || "someting went wrong";
    }

    return "someting went wrong";
  },
};
