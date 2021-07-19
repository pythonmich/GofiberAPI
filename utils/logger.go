package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

// TODO: Logger

// StandardLogger enforces specific log formats
type StandardLogger struct {
	*logrus.Logger
}

func NewLogger() *StandardLogger {
	baseLogger := logrus.New()
	standardLogger := &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.TextFormatter{}
	standardLogger.SetOutput(os.Stdout)
	standardLogger.SetLevel(logrus.DebugLevel)
	return standardLogger
}
