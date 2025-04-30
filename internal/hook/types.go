package hook

type IHook interface {
	Execute() error
	GetName() string
}
