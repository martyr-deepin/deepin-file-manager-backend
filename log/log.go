package log

import (
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("service/file-manager-backend")

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(tmpl string, args ...interface{}) {
	logger.Debugf(tmpl, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(tmpl string, args ...interface{}) {
	logger.Infof(tmpl, args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Warningf(tmpl string, args ...interface{}) {
	logger.Warningf(tmpl, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(tmpl string, args ...interface{}) {
	logger.Errorf(tmpl, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(tmpl string, args ...interface{}) {
	logger.Fatalf(tmpl, args...)
}
