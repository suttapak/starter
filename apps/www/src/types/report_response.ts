import { CommonModel } from "./common";

export interface ReportResponse extends CommonModel {
  code: string;
  name: string;
  display_name: string;
  icon: string;
  report_json_schema_type_id: number;
  report_json_schema_type: ReportJsonSchemaType;
}

export interface ReportJsonSchemaType extends CommonModel {
  name: string;
}
