package notification

type Level string

const (
	LevelInfo    Level = "info"
	LevelSuccess Level = "success"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

type Event struct {
	ID         string `json:"id"`
	ServiceID  string `json:"service_id"`
	Message    string `json:"message"`
	Level      Level  `json:"level"`
	OccurredAt int64  `json:"occurred_at"`
}

type Publisher interface {
	Publish(Event)
}
