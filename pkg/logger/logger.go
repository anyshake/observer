package logger

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetLogger(x any) *logrus.Entry {
	if v, ok := x.(string); ok {
		return logrus.WithFields(logrus.Fields{
			"module": strings.ToLower(v),
		})
	}

	val := reflect.ValueOf(x)
	if val.Kind() == reflect.Func {
		runtimeFunc := runtime.FuncForPC(val.Pointer())
		if runtimeFunc != nil {
			moduleNames := strings.Split(runtimeFunc.Name(), ".")
			if len(moduleNames) > 1 {
				lastPart := moduleNames[len(moduleNames)-1]
				moduleName := strings.Split(lastPart, "/")
				if len(moduleName) > 0 {
					return logrus.WithFields(logrus.Fields{
						"module": strings.ToLower(moduleName[len(moduleName)-1]),
					})
				}
			}
		}
	}

	return logrus.WithFields(logrus.Fields{
		"module": "unknown",
	})
}
