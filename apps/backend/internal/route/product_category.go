package route

import (
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"

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
