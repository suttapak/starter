package repository

import (
	"context"
	"strconv"

	"github.com/casbin/casbin/v2"
)

type (
	Permission interface {
		AddPolicy(ctx context.Context, groubId uint, routePath, method string) error
	}

	permission struct {
		enforcer *casbin.Enforcer
	}
)

// AddPolicy implements Permission.
func (p permission) AddPolicy(ctx context.Context, groubId uint, routePath string, method string) error {
	_, err := p.enforcer.AddPolicy(strconv.Itoa(int(groubId)), routePath, method)

	return err
}

func newPermissionRepository(enforcer *casbin.Enforcer) Permission {
	return permission{enforcer}
}
