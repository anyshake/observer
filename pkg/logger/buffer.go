package logger

import (
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/sirupsen/logrus"
)

type bufferHookImpl struct {
	formatter logrus.Formatter
	buffer    *ringbuf.Buffer[string]
}

func (h *bufferHookImpl) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *bufferHookImpl) Fire(entry *logrus.Entry) error {
	clone := entry.WithContext(entry.Context)
	clone.Message = entry.Message
	clone.Data = entry.Data
	clone.Level = entry.Level
	clone.Time = entry.Time

	b, err := h.formatter.Format(clone)
	if err != nil {
		return err
	}

	h.buffer.Push(string(b))
	return nil
}

func RegisterBufferLogger(bufSize int) *ringbuf.Buffer[string] {
	buf := ringbuf.New[string](bufSize)
	bufLogger := bufferHookImpl{
		formatter: &logrus.JSONFormatter{
			TimestampFormat: TIMESTAMP_FORMAT,
		},
		buffer: &buf,
	}

	logrus.AddHook(&bufLogger)
	return &buf
}
