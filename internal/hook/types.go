package hook

import "context"

type IHook interface {
	Execute(ctx context.Context) error
	GetName() string
}
