package logger

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type Formatter struct {
	log.TextFormatter
}

func (c *Formatter) Format(entry *log.Entry) ([]byte, error) {
	var logBuffer bytes.Buffer
	if entry.Buffer != nil {
		logBuffer = *entry.Buffer
	}

	logBuffer.WriteString(
		fmt.Sprintf(
			"%-7s [%s]",
			strings.ToUpper(entry.Level.String()),
			entry.Time.Format(c.TimestampFormat),
		),
	)

	if entry.HasCaller() {
		logBuffer.WriteString(fmt.Sprintf("[%s:%d]", entry.Caller.Function, entry.Caller.Line))
	}
	logBuffer.WriteString(fmt.Sprintf(" %s\n", entry.Message))
	return logBuffer.Bytes(), nil
}

func Init() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		_ = os.Mkdir("logs", os.ModePerm)
	}

	//logFile := &lumberjack.Logger{
	//	Filename:   "logs/app.log",
	//	MaxSize:    100,
	//	MaxBackups: 3,
	//	MaxAge:     28,
	//	Compress:   true,
	//}

	//log.SetOutput(logFile)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&Formatter{
		log.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
			ForceColors:     false,
			DisableColors:   true,
		},
	})
}
