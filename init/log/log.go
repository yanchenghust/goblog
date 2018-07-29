package log

import "github.com/yanchenghust/goblog/log"

var logger *log.Log

func InitLog() {
	fp := log.NewFileProvider("/www/log/blog")
	logger = log.NewLog(log.LvlDebug, fp)
	cl := log.NewConsoleProvider()
	logger.AppendProvider(cl)
	logger.Start()
}

func StopLog() {
	logger.Stop()
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
