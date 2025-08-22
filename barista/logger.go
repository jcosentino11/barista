package barista

import (
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type ConsoleLogger struct {
	*log.Logger
}

func NewConsoleLogger(name string) *ConsoleLogger {
	return &ConsoleLogger{
		log.New(os.Stderr, name+": ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
