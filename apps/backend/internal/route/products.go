package route

import (
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

type (
	products struct {
		r        *gin.Engine
		guard    middleware.AuthGuardMiddleware
		products controller.Products
	}
)

func newProducts(r *gin.Engine, productsController controller.Products, guard middleware.AuthGuardMiddleware) *products {
	return &products{
		r:        r,
		products: productsController,
		guard:    guard,
	}
}

func useProducts(a *products) {
	r := a.r.Group("teams/:team_id/products", a.guard.Protect, a.guard.TeamPermission)
	{
		r.GET("/:products_id", a.products.GetProduct)
		r.GET("", a.products.GetProducts)
		r.POST("", a.products.CreateProducts)
		r.PUT("/:products_id", a.products.UpdateProducts)
		r.POST("/:products_id/upload_image", a.products.UploadProductImages)
		r.DELETE("/:products_id", a.products.DeleteProducts)
		r.DELETE("/:products_id/product_image/:product_image_id", a.products.DeleteProductImage)
	}
}
