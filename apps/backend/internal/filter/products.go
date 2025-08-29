package filter

type (
	ProductsFilter struct {
		Code string `form:"code"`
		Name string `form:"name"`
		UOM  string `form:"uom"`
	}
)
