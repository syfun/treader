package logging

import (
	"github.com/sirupsen/logrus"
)

// Init logger.
func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
	})
}

// GetLogger get logger with topic.
func GetLogger(topic string) *logrus.Entry {
	return logrus.WithField("topic", topic)
}
