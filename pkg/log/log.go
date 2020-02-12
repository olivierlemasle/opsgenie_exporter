// Heavily inspired by "github.com/prometheus/common/log"

package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var logger = logrus.New()

type loggerSettings struct {
	level  string
	format string
}

func (s *loggerSettings) apply(ctx *kingpin.ParseContext) error {
	level, err := logrus.ParseLevel(s.level)
	if err != nil {
		return err
	}
	logger.SetLevel(level)

	switch s.format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "txt":
	default:
		return fmt.Errorf("Invalid log format %v", s.format)
	}
	return nil
}

// AddFlags adds the log flags used by this package to the Kingpin application.
// To use the default Kingpin application, call AddFlags(kingpin.CommandLine)
func AddFlags(a *kingpin.Application) {
	s := loggerSettings{}
	a.Flag("log.level", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]").
		Default("info").
		StringVar(&s.level)
	a.Flag("log.format", "Output format of log messages. One of: [txt, json]").
		Default("txt").
		StringVar(&s.format)
	a.Action(s.apply)
}

// Logger returns the configured logger
func Logger() *logrus.Logger {
	return logger
}
