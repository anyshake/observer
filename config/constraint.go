package config

import "github.com/anyshake/observer/internal/dao/action"

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

	Init(handler *action.Handler) error // should be initialized during startup hook
	Set(handler *action.Handler, newVal any) error
	Get(handler *action.Handler) (any, error)
	Restore(handler *action.Handler) error
}
