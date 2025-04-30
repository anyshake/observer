package graph_resolver

import (
	"context"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/seisevent"
	"github.com/anyshake/observer/pkg/timesource"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ContextKey string

type Resolver struct {
	HardwareDev              hardware.IHardware
	TimeSource               *timesource.Source
	ActionHandler            *action.Handler
	StationConfigConstraints []config.IConstraint
	ServiceMap               map[string]service.IService
	SeisEventSource          map[string]seisevent.IDataSource
}

func (r *Resolver) getCurrentUserId(ctx context.Context) string {
	user := ctx.Value(ContextKey("user_status"))
	if user == nil {
		return ""
	}

	userStatusMap, ok := user.(map[string]any)
	if !ok {
		return ""
	}

	user, ok = userStatusMap["user_id"]
	if !ok {
		return ""
	}

	return user.(string)
}

func (r *Resolver) checkIsAdmin(ctx context.Context) bool {
	user := ctx.Value(ContextKey("user_status"))
	if user == nil {
		return false
	}

	userStatusMap, ok := user.(map[string]any)
	if !ok {
		return false
	}

	isAdmin, ok := userStatusMap["is_admin"]
	if !ok {
		return false
	}

	return isAdmin.(bool)
}
