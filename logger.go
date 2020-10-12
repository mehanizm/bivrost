package bivrost

import (
	"fmt"
	"log"
	"os"
)

// Logger is implemented by any logging system that is used for standard logs.
type Logger interface {
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

type LoggingLevel int

const (
	DEBUG LoggingLevel = iota
	INFO
	WARNING
	ERROR
)

type DefaultLog struct {
	*log.Logger
	level LoggingLevel
}

func DefaultLogger(level LoggingLevel, componentName string) *DefaultLog {
	return &DefaultLog{Logger: log.New(os.Stderr, fmt.Sprintf("bivrost [%s] ", componentName), log.Lmsgprefix|log.LstdFlags), level: level}
}

func (l *DefaultLog) Errorf(f string, v ...interface{}) {
	if l.level <= ERROR {
		l.Printf("ERROR: "+f, v...)
	}
}

func (l *DefaultLog) Warningf(f string, v ...interface{}) {
	if l.level <= WARNING {
		l.Printf("WARNING: "+f, v...)
	}
}

func (l *DefaultLog) Infof(f string, v ...interface{}) {
	if l.level <= INFO {
		l.Printf("INFO: "+f, v...)
	}
}

func (l *DefaultLog) Debugf(f string, v ...interface{}) {
	if l.level <= DEBUG {
		l.Printf("DEBUG: "+f, v...)
	}
}
