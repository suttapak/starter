package seed

import (
	"context"

	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

func SeedRole(db *gorm.DB) error {

	ctx := context.Background()
	count, err := gorm.G[model.Role](db).Count(ctx, "id")
	if err != nil {
		return err
	}
	if count == 0 {
		var roles = []model.Role{{
			CommonModel: model.CommonModel{ID: idx.RoleUser},
			Name:        "User",
		}, {
			CommonModel: model.CommonModel{ID: idx.RoleModerator},
			Name:        "Moderator",
		}, {
			CommonModel: model.CommonModel{ID: idx.RoleAdmin},
			Name:        "Admin",
		}, {
			CommonModel: model.CommonModel{ID: idx.RoleSuperAdmin},
			Name:        "SuperAdmin",
		}}
		if err := gorm.G[model.Role](db).CreateInBatches(ctx, &roles, len(roles)); err != nil {
			return err
		}
	}
	return nil
}
