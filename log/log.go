package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Level uint32

const (
	FieldLogID           = "log_id"
	FieldEndpoint        = "endpoint"
	FieldMethod          = "method"
	FieldRequestBody     = "request_body"
	FieldRequestHeaders  = "request_headers"
	FieldResponseBody    = "response_body"
	FieldResponseHeaders = "response_headers"
)

type message struct {
	Message  interface{} `json:"message"`
	Level    Level       `json:"level"`
	File     string      `json:"file"`
	FuncName string      `json:"func"`
	Line     int         `json:"line"`
}

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

func parseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

func (level *Level) UnmarshalText(text []byte) error {
	l, err := parseLevel(string(text))
	if err != nil {
		return err
	}

	*level = Level(l)

	return nil
}

func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case TraceLevel:
		return []byte("trace"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid logrus level %d", level)
}

type Logger struct {
	logger *log.Logger
	field  Field
	fields log.Fields
}

type Field struct {
	// LogID, should be unique id
	LogID          string
	Endpoint       string
	Method         string
	RequestBody    interface{}
	RequestHeader  interface{}
	ResponseBody   interface{}
	ResponseHeader interface{}
}

func (l *Logger) SetLogID(logID string) *Logger {
	l.field.LogID = logID
	return l
}

func (l *Logger) SetEndpoint(endpoint string) *Logger {
	l.field.Endpoint = endpoint
	return l
}

func (l *Logger) SetMethod(method string) *Logger {
	l.field.Method = method
	return l
}

func (l *Logger) SetRequestBody(body interface{}) *Logger {
	l.field.RequestBody = body
	return l
}

func (l *Logger) SetRequestHeaders(headers interface{}) *Logger {
	l.field.RequestHeader = headers
	return l
}

func (l *Logger) SetResponseBody(body interface{}) *Logger {
	l.field.ResponseBody = body
	return l
}

func (l *Logger) SetResponseHeaders(headers interface{}) *Logger {
	l.field.ResponseHeader = headers
	return l
}

func (l *Logger) AddMessage(level Level, message interface{}) *Logger {
	l.setCaller(message, level)
	return l
}

func (l *Logger) Print() {
	fill := fillField(l.fields)
	fill(FieldLogID, l.field.LogID)
	fill(FieldEndpoint, l.field.Endpoint)
	fill(FieldMethod, l.field.Method)
	fill(FieldRequestBody, l.field.RequestBody)
	fill(FieldRequestHeaders, l.field.RequestHeader)
	fill(FieldResponseBody, l.field.ResponseBody)
	fill(FieldResponseHeaders, l.field.ResponseHeader)

	entry := l.logger.WithFields(l.fields)

	messages := ensureStackType(l.fields["stack"])
	entry.Tracef("%+v", messages[0].Message)
	l.field = Field{}
	return
}

func (l *Logger) setCaller(msg interface{}, level Level) {
	if msg != nil && msg != "" {
		if pc, file, line, ok := runtime.Caller(2); ok {
			fName := runtime.FuncForPC(pc).Name()

			err, ok := msg.(error)
			if ok && err != nil {
				msg = err.Error()
			}

			msg := message{
				Message:  msg,
				Level:    level,
				File:     file,
				FuncName: fName,
				Line:     line,
			}

			l.addMessageStack(msg)
		}
	}
}

func (l *Logger) addMessageStack(msg message) {
	stack := ensureStackType(l.fields["stack"])
	l.fields["stack"] = append(stack, msg)
}

func fillField(m map[string]interface{}) func(string, interface{}) {
	return func(key string, val interface{}) {
		if val == nil {
			return
		}

		m[key] = val
	}
}

func ensureStackType(stack interface{}) (val []message) {
	val, ok := stack.([]message)
	if !ok {
		panic("perkakas/log: stack trace is expecting a message but found other types")
	}

	return
}

func newLog(formatter log.Formatter, out io.Writer, level log.Level, reportCaller bool) (l *log.Logger) {
	l = log.New()
	l.SetFormatter(formatter)
	l.SetOutput(out)
	l.SetLevel(level)

	return
}

func NewLogger() (logger *Logger) {
	formatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
		// PrettyPrint:     true,
		FieldMap: log.FieldMap{
			log.FieldKeyMsg: "log_message",
		},
		DisableTimestamp: true,
	}

	newLogger := newLog(formatter, os.Stdout, log.TraceLevel, false)

	logger = new(Logger)
	logger.logger = newLogger
	logger.fields = make(map[string]interface{})
	logger.fields["stack"] = []message{}

	return
}
