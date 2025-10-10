package model

type (
	Product struct {
		CommonModel
		TeamID                 uint                     `db:"team_id" json:"team_id"`
		Code                   string                   `db:"code" json:"code"`
		Name                   string                   `db:"name" json:"name"`
		Description            string                   `db:"description" json:"description"`
		UOM                    string                   `db:"uom" json:"uom"`
		Price                  int64                    `db:"price" json:"price"` // Price in cents
		ProductProductCategory []ProductProductCategory `db:"-" json:"product_product_category,omitempty"`
		ProductImage           []ProductImage           `db:"-" json:"product_image,omitempty"`
	}

	ProductImage struct {
		CommonModel
		ProductID uint   `db:"product_id" json:"product_id"`
		ImageID   uint   `db:"image_id" json:"image_id"`
		Image     *Image `db:"-" json:"image,omitempty"`
	}
)
