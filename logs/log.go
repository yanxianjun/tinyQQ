package logs

import "github.com/sirupsen/logrus"


var Logger = logrus.WithField("bot", "service")

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	//logrus.SetLevel(logrus.DebugLevel)
}