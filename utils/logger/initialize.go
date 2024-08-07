package logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func Initialize() {
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: TIMESTAMP_FORMAT,
		FieldsOrder:     []string{"module"},
	})
}
