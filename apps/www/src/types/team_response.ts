import { CommonModel } from "./common";

export interface TeamResponse extends CommonModel {
  name: string;
  username: string;
  email: string;
  phone: string;
  address: string;
  description: string;
}
