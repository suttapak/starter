import { CommonModel } from "./common";

export interface Image extends CommonModel {
  path: string;
  url: string;
  size: number;
  width: number;
  height: number;
  type: string;
}
