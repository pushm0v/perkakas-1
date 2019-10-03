package log

import (
	"io"
	"os"
	"runtime"
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
	FieldMessage         = "custom_message"
	FieldLevel           = "message_level"
)

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

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

	// Log level when print
	Level   Level
	Message interface{}
}

// Set default logger level
func (l *Logger) SetLoggerLevel(lv Level) {
	l.logger.Level = log.Level(uint32(lv))
}

func (l *Logger) Set(field Field) *Logger {
	l.field = field
	l.setCaller(l.field.Message)

	return l
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

func (l *Logger) SetMessage(level Level, message interface{}) *Logger {
	l.field.Message = message
	l.field.Level = level
	l.setCaller(l.field.Message)

	return l
}

// func (l *Logger) Log(lv Level, args ...interface{}) {
func (l *Logger) Print(args ...interface{}) {
	l.fields[FieldLogID] = l.field.LogID
	l.fields[FieldEndpoint] = l.field.Endpoint
	l.fields[FieldMethod] = l.field.Method
	l.fields[FieldRequestBody] = l.field.RequestBody
	l.fields[FieldRequestHeaders] = l.field.RequestHeader
	l.fields[FieldResponseBody] = l.field.ResponseBody
	l.fields[FieldResponseHeaders] = l.field.ResponseHeader
	l.fields[FieldLevel] = l.field.Level
	l.fields[FieldMessage] = l.field.Message

	entry := l.logger.WithFields(l.fields)

	if l.field.Level > WarnLevel {
		entry.Logger.SetOutput(os.Stdout)
	} else {
		entry.Logger.SetOutput(os.Stderr)
	}

	entry.Logln(log.Level(uint32(l.field.Level)), args...)
	l.field = Field{}
	return
}

func (l *Logger) setCaller(errorMessage interface{}) {
	if errorMessage != nil && errorMessage != "" {
		if pc, file, line, ok := runtime.Caller(2); ok {
			fName := runtime.FuncForPC(pc).Name()
			l.fields["file"] = file
			l.fields["line"] = line
			l.fields["func"] = fName
		}
	}
}

func newLog(formatter log.Formatter, out io.Writer, level log.Level, reportCaller bool) (l *log.Logger) {
	l = log.New()
	l.SetFormatter(formatter)
	l.SetOutput(out)
	l.SetLevel(level)
	l.SetReportCaller(reportCaller)

	return
}

func NewLogger() (logger *Logger) {
	formatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: log.FieldMap{
			log.FieldKeyMsg: "log_message",
		},
		DisableTimestamp: true,
	}

	newLogger := newLog(formatter, os.Stdout, log.TraceLevel, false)

	logger = new(Logger)
	logger.logger = newLogger
	logger.fields = make(map[string]interface{})

	return
}
