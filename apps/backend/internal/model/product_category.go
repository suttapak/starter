package model

type (
	ProductProductCategory struct {
		CommonModel
		ProductID         uint             `db:"product_id" json:"product_id"`
		ProductCategoryID uint             `db:"product_category_id" json:"category_id"`
		ProductCategory   *ProductCategory `db:"-" json:"category,omitempty"`
	}

	ProductCategory struct {
		CommonModel
		TeamID uint   `db:"team_id" json:"team_id"`
		Name   string `db:"name" json:"name"`
	}
)
