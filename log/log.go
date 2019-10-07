package log

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/kitabisa/perkakas/v2/httputil"
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
	logger           *log.Logger
	fields           log.Fields
	autoGenRequestID bool
}

func (l *Logger) SetRequest(req interface{}) {
	l.fields[FieldLogID] = uuid.NewV1().String()

	switch v := req.(type) {
	case *http.Request:
		l.fields[FieldEndpoint] = v.URL.String()
		l.fields[FieldMethod] = v.Method
		l.fields[FieldRequestHeaders] = v.Header

		switch v.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			l.fields[FieldRequestBody] = httputil.ReadRequestBody(v)
		}
	default:
		l.fields[FieldRequestBody] = req
	}
}

func (l *Logger) SetResponse(res interface{}, body interface{}) {
	switch v := res.(type) {
	case *http.Response:
		l.fields[FieldResponseHeaders] = v.Header
		l.fields[FieldResponseBody] = body
	}
}

func (l *Logger) AddMessage(level Level, message interface{}) *Logger {
	l.setCaller(message, level)
	return l
}

func (l *Logger) Print() {
	messages := ensureStackType(l.fields["stack"])
	if len(messages) == 0 || len(l.fields) == 0 {
		return
	}

	entry := l.logger.WithFields(l.fields)
	entry.Tracef("%+v", messages[0].Message)
	l.fields = make(map[string]interface{})
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
