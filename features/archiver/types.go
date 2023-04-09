package archiver

import "com.geophone.observer/features/collector"

type ArchiverOptions struct {
	Enable             bool
	Database           int
	Password           string
	OnCompleteCallback func()
	OnErrorCallback    func(error)
	Message            *collector.Message
}
