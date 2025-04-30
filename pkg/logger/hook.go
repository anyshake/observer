package logger

import "github.com/sirupsen/logrus"

func (hook *rotateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *rotateHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.logWriter.Write(b)
	return err
}
