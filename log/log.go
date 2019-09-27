package log

import (
	"context"
	"io"
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

type Level uint32

const (
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
	ctx    context.Context
	field  Field
	fields log.Fields
}

type Field struct {
	Endpoint       string
	Method         string
	RequestBody    interface{}
	RequestHeader  interface{}
	ResponseBody   interface{}
	ResponseHeader interface{}
	ErrorMessage   interface{}
}

func (l *Logger) Set(field Field) {
	l.field = field
}

func (l *Logger) SetLevel(lv Level) {
	l.logger.Level = log.Level(uint32(lv))
}

func (l *Logger) SetEndpoint(endpoint string) {
	l.field.Endpoint = endpoint
}

func (l *Logger) SetMethod(method string) {
	l.field.Method = method
}

func (l *Logger) SetRequestBody(body interface{}) {
	l.field.RequestBody = body
}

func (l *Logger) SetRequestHeaders(headers interface{}) {
	l.field.RequestHeader = headers
}

func (l *Logger) SetResponseBody(body interface{}) {
	l.field.ResponseBody = body
}

func (l *Logger) SetResponseHeaders(headers interface{}) {
	l.field.ResponseHeader = headers
}

func (l *Logger) SetErrorMessage(errorMessage interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		l.fields["file"] = file
		l.fields["line"] = line
		l.fields["func"] = fName
	}
	l.field.ErrorMessage = errorMessage
}

func (l *Logger) Log(lv Level, args ...interface{}) {
	l.fields = map[string]interface{}{
		"log_id":             l.ctx.Value("log_id"),
		FieldEndpoint:        l.field.Endpoint,
		FieldMethod:          l.field.Method,
		FieldRequestBody:     l.field.RequestBody,
		FieldRequestHeaders:  l.field.RequestHeader,
		FieldResponseBody:    l.field.ResponseBody,
		FieldResponseHeaders: l.field.ResponseHeader,
		FieldErrorMessage:    l.field.ErrorMessage,
	}

	if l.fields[FieldErrorMessage] == nil {
		if pc, file, line, ok := runtime.Caller(1); ok {
			fName := runtime.FuncForPC(pc).Name()
			l.fields["file"] = file
			l.fields["line"] = line
			l.fields["func"] = fName
		}
	}

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

func NewLogger(ctx context.Context) (logger *Logger) {
	formatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	newLogger := newLog(formatter, os.Stdout, log.TraceLevel, false)

	logger = new(Logger)
	logger.logger = newLogger
	logger.fields = make(map[string]interface{})
	logger.ctx = ctx

	return
}
