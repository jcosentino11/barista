package barista

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Verbosef(format string, v ...interface{})
}

type ConsoleLogger struct {
	Verbose bool
	*log.Logger
}

func NewConsoleLogger(name string) *ConsoleLogger {
	return &ConsoleLogger{
		false,
		log.New(os.Stderr, name+": ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *ConsoleLogger) Verbosef(format string, v ...interface{}) {
	if l.Verbose {
		l.Printf(format, v...)
	}
}
