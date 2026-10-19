package config

import (
	"context"

	"github.com/anyshake/observer/internal/dao/action"
)

type IConstraint interface {
	IsRequired() bool // double check in setter implementation
	GetDefaultValue() any
	GetNamespace() string
	GetDescription() string
	GetName() string

	GetType() action.SettingType
	GetKey() string // unique identifier
	GetVersion() int
	GetOptions() map[string]any

	Init(ctx context.Context, handler *action.Handler) error // should be initialized during startup hook
	Set(ctx context.Context, handler *action.Handler, newVal any) error
	Get(ctx context.Context, handler *action.Handler) (any, error)
	Restore(ctx context.Context, handler *action.Handler) error
}
