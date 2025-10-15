import type { CommonModel } from "@/shared/types";

export interface TeamResponse extends CommonModel {
  name: string;
  username: string;
  email: string;
  phone: string;
  address: string;
  description: string;
}
