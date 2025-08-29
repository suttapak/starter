package model

type (
	ProductProductCategory struct {
		CommonModel
		ProductID         uint            `json:"product_id"`
		ProductCategoryID uint            `json:"category_id"`
		ProductCategory   ProductCategory `json:"category"`
	}

	ProductCategory struct {
		CommonModel
		TeamID uint   `json:"team_id"`
		Name   string `json:"name"`
	}
)
