package model

type (
	Product struct {
		CommonModel
		TeamID                 uint                     `json:"team_id"`
		Code                   string                   `json:"code"`
		Name                   string                   `json:"name"`
		Description            string                   `json:"description"`
		UOM                    string                   `json:"uom"`
		Price                  int64                    `json:"price" gorm:"comment: 'Price in cents'"`
		ProductProductCategory []ProductProductCategory `json:"product_product_category"`
		ProductImage           []ProductImage           `json:"product_image"`
	}

	ProductImage struct {
		CommonModel
		ProductID uint  `json:"product_id"`
		ImageID   uint  `json:"image_id"`
		Image     Image `json:"image"`
	}
)
