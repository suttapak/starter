package route

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

type (
	productCategory struct {
		r               *gin.Engine
		guard           middleware.AuthGuardMiddleware
		productCategory controller.ProductCategory
	}
)

func newProductCategory(r *gin.Engine, productCategoryController controller.ProductCategory, guard middleware.AuthGuardMiddleware) *productCategory {
	return &productCategory{
		r:               r,
		productCategory: productCategoryController,
		guard:           guard,
	}
}

func useProductCategory(a *productCategory) {
	r := a.r.Group("teams/:team_id/product_category", a.guard.Protect, a.guard.TeamPermission)
	{
		r.GET("/:product_category_id", a.productCategory.GetProductCategory)
		r.GET("", a.productCategory.GetProductCategories)
		r.POST("", a.productCategory.CreateProductCategory)
		r.PUT("/:product_category_id", a.productCategory.UpdateProductCategory)
		r.DELETE("/:product_category_id", a.productCategory.DeleteProductCategory)
	}
}

func seedProductCategoryPermission(db *gorm.DB) {
	/*
		TeamRoleOwnerID = iota + 1
		TeamRoleAdminID
		TeamRoleMemberID
	*/
	// v1 for role id, v2 for path, v3 for method
	var permission = []gormadapter.CasbinRule{
		{
			Ptype: "p",
			V0:    "3",
			V1:    "/teams/{id}/product_category/*",
			V2:    "GET",
		},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&permission)

}
