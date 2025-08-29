package filter

type (
	TeamMemberFilter struct {
		Username string `json:"username"`
	}
	TeamFilter struct {
		Name string `form:"name"`
	}
	PublicProductFilter struct {
		Code string `form:"code"`
		Name string `form:"name"`
		UOM  string `form:"uom"`
	}
)
