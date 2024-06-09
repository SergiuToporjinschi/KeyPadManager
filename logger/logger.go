package logger

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

// CustomTextFormatter is a custom formatter that displays only values without attribute names
type CustomTextFormatter struct {
	TimestampFormat string
}

// Format formats the log entry without attribute names
func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := ""
	if !entry.HasCaller() {
		timestamp = entry.Time.Format(f.TimestampFormat)
	}
	// Extract log level
	level := strings.ToUpper(entry.Level.String())
	return []byte(timestamp + " " + level + " " + entry.Message + "\n"), nil
}

func Init() {
	Log = logrus.New()

	// Set the formatter, output, and level for the logger
	// Log.SetFormatter(&logrus.TextFormatter{
	// 	FullTimestamp: true,
	// })
	// Set the formatter, output, and level for the logger
	Log.SetFormatter(&CustomTextFormatter{
		TimestampFormat: "2006-01-02T15:04:05", // Customize timestamp format
	})

	// Define file for logging
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Fatal("Unable to open log file:", err)
	}

	// Set the output to both console and file
	Log.SetOutput(io.MultiWriter(os.Stdout, file))

	// Set the log level
	Log.SetLevel(logrus.DebugLevel)
	Log.Infof("Logger initialized %s", "successfully")
	log.SetOutput(Log.Writer())
}
