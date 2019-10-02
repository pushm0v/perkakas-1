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
	FieldErrorMessage    = "error_message"
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
	LogID          string
	Endpoint       string
	Method         string
	RequestBody    interface{}
	RequestHeader  interface{}
	ResponseBody   interface{}
	ResponseHeader interface{}
	ErrorMessage   interface{}
}

func (l *Logger) SetLevel(lv Level) {
	l.logger.Level = log.Level(uint32(lv))
}

func (l *Logger) Set(field Field) *Logger {
	l.field = field
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

func (l *Logger) SetErrorMessage(errorMessage interface{}) *Logger {
	l.field.ErrorMessage = errorMessage
	return l
}

func (l *Logger) Trace(args ...interface{}) {
	l.log(TraceLevel, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(InfoLevel, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log(WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(FatalLevel, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log(PanicLevel, args...)
}

func (l *Logger) log(lv Level, args ...interface{}) {
	l.fields = map[string]interface{}{
		FieldLogID:           l.field.LogID,
		FieldEndpoint:        l.field.Endpoint,
		FieldMethod:          l.field.Method,
		FieldRequestBody:     l.field.RequestBody,
		FieldRequestHeaders:  l.field.RequestHeader,
		FieldResponseBody:    l.field.ResponseBody,
		FieldResponseHeaders: l.field.ResponseHeader,
		FieldErrorMessage:    l.field.ErrorMessage,
	}

	setCaller(l.fields)

	entry := l.logger.WithFields(l.fields)

	if lv > WarnLevel {
		entry.Logger.SetOutput(os.Stdout)
	} else {
		entry.Logger.SetOutput(os.Stderr)
	}

	entry.Logln(log.Level(uint32(lv)), args...)
	l.field = Field{}
	return
}

func newLog(formatter log.Formatter, out io.Writer, level log.Level, reportCaller bool) (l *log.Logger) {
	l = log.New()
	l.SetFormatter(formatter)
	l.SetOutput(out)
	l.SetLevel(level)
	l.SetReportCaller(reportCaller)

	return
}

func setCaller(fields map[string]interface{}) {
	if fields[FieldErrorMessage] != nil && fields[FieldErrorMessage] != "" {
		if pc, file, line, ok := runtime.Caller(3); ok {
			fName := runtime.FuncForPC(pc).Name()
			fields["file"] = file
			fields["line"] = line
			fields["func"] = fName
		}
	}
}

func NewLogger() (logger *Logger) {
	formatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	newLogger := newLog(formatter, os.Stdout, log.TraceLevel, false)

	logger = new(Logger)
	logger.logger = newLogger
	logger.fields = make(map[string]interface{})

	return
}
