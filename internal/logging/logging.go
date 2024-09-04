package logging

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

var Log *logrus.Logger

func Init() error {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{})

	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "ukip")
	if err != nil {
		return err
	}

	Log.AddHook(hook)
	return nil
}