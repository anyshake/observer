package logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func Init() {
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: TIMESTAMP_FORMAT,
		FieldsOrder:     []string{"module"},
	})
}
