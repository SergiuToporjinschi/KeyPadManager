package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strconv"
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

func NewSLogger() *slog.Logger {
	return slog.New(
		NewPrettyConsoleHandler(*HandleCustomLevels()),
	)
}

const (
	timeFormat  = "[2006-01-02 15:04:05.000]"
	red         = 31
	green       = 32
	lightRed    = 91
	lightGreen  = 92
	lightYellow = 93
	lightBlue   = 94
	white       = 97
)

type PrettyConsoleHandler struct {
	slog.Handler
	log       *log.Logger
	useColors bool
	addSource bool
	options   *slog.HandlerOptions
}

func NewPrettyConsoleHandler(opts slog.HandlerOptions) *PrettyConsoleHandler {
	return &PrettyConsoleHandler{
		options:   &opts,
		useColors: true,
		addSource: true,
		log:       log.New(os.Stdout, "", 0),
	}
}
func (h *PrettyConsoleHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.options.Level.Level()
}

func (h *PrettyConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyConsoleHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *PrettyConsoleHandler) WithGroup(name string) slog.Handler {
	return &PrettyConsoleHandler{Handler: h.Handler.WithGroup(name)}
}

func (h *PrettyConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	level := ReplaceLevelName(r.Level) + ":"
	buffer := make([]byte, 0, 10)
	appendToBuffer := func(content string) {
		buffer = append(buffer, []byte(content)...)
		buffer = append(buffer, ' ')
	}

	appendToBuffer(r.Time.Format(timeFormat))

	if h.useColors {
		switch r.Level {
		case slog.LevelDebug:
			level = colorize(lightGreen, level)
		case slog.LevelInfo:
			level = colorize(lightBlue, level)
		case slog.LevelWarn:
			level = colorize(lightYellow, level)
		case slog.LevelError:
			level = colorize(red, level)
		}
	}

	appendToBuffer(level)

	if h.addSource {
		appendToBuffer(h.getSourceString(r.PC))
	}

	if h.useColors {
		appendToBuffer(colorize(white, r.Message))
	} else {
		appendToBuffer(r.Message)
	}

	errorStack := ""
	attrs := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		switch a.Value.Any().(type) {
		case error:
			attrs[a.Key] = a.Value.Any().(error).Error()
			errorStack = h.getStack()
		default:
			if reflect.ValueOf(a.Value.Any()).Kind() == reflect.Struct {
				obj, err := json.Marshal(a.Value.Any())
				if err == nil {
					attrs[a.Key] = string(obj)
				} else {
					attrs[a.Key] = a.Value.Any()
				}
			} else {
				attrs[a.Key] = a.Value.Any()
			}
		}
		return true
	})

	appendToBuffer(h.getAttributesString(attrs))

	if (r.Level == slog.LevelError) && errorStack != "" {
		appendToBuffer("\n" + errorStack)
	}

	h.log.Println(string(buffer))
	return nil
}

func (h *PrettyConsoleHandler) getSourceString(caller uintptr) string {
	fs := runtime.CallersFrames([]uintptr{caller})
	f, _ := fs.Next()
	shortFileName := f.File[strings.LastIndex(f.File, "/")+1:]
	return fmt.Sprintf("%s:%d", shortFileName, f.Line)
}

func (h *PrettyConsoleHandler) getAttributesString(fields map[string]interface{}) string {
	var attributeList []string
	for key, value := range fields {
		if key == "source" || key == "level" || key == "msg" || key == "time" {
			continue
		}
		attributeList = append(attributeList, fmt.Sprintf("%s: %v", key, value))
	}
	var attrStr string
	if len(attributeList) != 0 {
		attrStr = fmt.Sprintf("[%s]", strings.Join(attributeList, ", "))
	}
	return attrStr
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, "\033[0m")
}

func (h *PrettyConsoleHandler) getStack() string {
	var stackTrace string
	pcs := make([]uintptr, 15)
	n := runtime.Callers(8, pcs) // Skip first 8 frames, are referring to logger stack trace
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		stackTrace += fmt.Sprintf("%s:%d\n", frame.File, frame.Line)
		if !more {
			break
		}
	}

	return stackTrace

}

var LevelNames = map[slog.Level]string{
	slog.LevelDebug: "DEBUG",
	slog.LevelInfo:  "INFO",
	slog.LevelWarn:  "WARNING",
	slog.LevelError: "ERROR",
}

func ReplaceLevelName(level slog.Level) string {
	levelLabel, exists := LevelNames[level]
	if !exists {
		levelLabel = level.String()
	}
	return levelLabel
}
func HandleCustomLevels() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource: true,
		Level:     findLeveler("DEBUG"),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				a.Value = slog.StringValue(ReplaceLevelName(a.Value.Any().(slog.Level)))
			}
			return a
		},
	}
}
func findLeveler(strLevel string) slog.Level {
	for key := range LevelNames {
		if LevelNames[key] == strLevel {
			return key
		}
	}
	return slog.LevelError
}
