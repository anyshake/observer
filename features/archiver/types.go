package archiver

import "com.geophone.observer/features/collector"

type ArchiverOptions struct {
	Enable             bool
	OnCompleteCallback func()
	OnErrorCallback    func(error)
	Message            *collector.Message
}
