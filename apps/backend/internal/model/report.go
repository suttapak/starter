package model

import "github.com/suttapak/starter/internal/idx"

type (
	ReportTemplate struct {
		CommonModel
		Code                   string                   `db:"code" json:"code"`
		Name                   string                   `db:"name" json:"name"`
		DisplayName            string                   `db:"display_name" json:"display_name"`
		Icon                   string                   `db:"icon" json:"icon"`
		ReportJsonSchemaTypeID idx.ReportJsonSchemaType `db:"report_json_schema_type_id" json:"report_json_schema_type_id"`
		ReportJsonSchemaType   *ReportJsonSchemaType    `db:"-" json:"report_json_schema_type,omitempty"`
	}

	ReportJsonSchemaType struct {
		CommonModel
		Name string `db:"name" json:"name"`
	}
)
