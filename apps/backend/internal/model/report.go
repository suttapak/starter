package model

import "github.com/suttapak/starter/internal/idx"

type (
	ReportTemplate struct {
		CommonModel
		Code                   string                   `json:"code"`
		Name                   string                   `json:"name"`
		DisplayName            string                   `json:"display_name"`
		Icon                   string                   `json:"icon"`
		ReportJsonSchemaTypeID idx.ReportJsonSchemaType `json:"report_json_schema_type_id"`
		ReportJsonSchemaType   ReportJsonSchemaType     `json:"report_json_schema_type"`
	}

	ReportJsonSchemaType struct {
		CommonModel
		Name string `json:"name"`
	}
)
