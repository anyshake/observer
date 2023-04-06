package archiver

import "com.geophone.observer/features/collector"

type ArchiverOptions struct {
	Enable             bool
	Name               string
	OnCompleteCallback func()
	OnErrorCallback    func(error)
	Status             *collector.Status
	Message            *collector.Message
}
