package log

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/kitabisa/perkakas/v2/httputil"
	"github.com/kitabisa/perkakas/v2/token/jwt"
)

type Level uint32

const (
	FieldLogID           = "log_id"
	FieldHTTPStatus      = "http_status"
	FieldEndpoint        = "endpoint"
	FieldMethod          = "method"
	FieldServiceName     = "service"
	FieldUserID          = "user_id"
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
	logger *log.Logger
	fields sync.Map
	id     string
}

func (l *Logger) NewChildLogger() (logger *Logger) {
	logger = newLogger(l.id)
	return
}

func (l *Logger) SetRequest(req interface{}) {
	switch v := req.(type) {
	case *http.Request:
		token, ok := v.Context().Value("token").(*jwt.UserClaim)
		if ok {
			l.fields.Store(FieldUserID, token.UserID)
		}

		l.fields.Store(FieldEndpoint, v.URL.String())
		l.fields.Store(FieldMethod, v.Method)
		l.fields.Store(FieldRequestHeaders, v.Header)

		switch v.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			l.fields.Store(FieldRequestBody, httputil.ReadRequestBody(v))
		}
	default:
		l.fields.Store(FieldRequestBody, req)
	}
}

func (l *Logger) SetResponse(res interface{}, body []byte) {
	switch v := res.(type) {
	case http.ResponseWriter:
		l.fields.Store(FieldResponseHeaders, v.Header())
		l.fields.Store(FieldResponseBody, string(body))
	case *http.Response:
		l.fields.Store(FieldResponseHeaders, v.Header)
		l.fields.Store(FieldResponseBody, string(body))
	}
}

func (l *Logger) AddMessage(level Level, message interface{}) *Logger {
	l.setCaller(message, level)
	return l
}

func (l *Logger) Print() {
	stackVal, _ := l.fields.Load("stack")
	messages := ensureStackType(stackVal)

	tempLoggerFields := make(map[string]interface{})
	fieldsLen := 0

	l.fields.Range(func(key interface{}, value interface{}) bool {
		tempLoggerFields[key.(string)] = value
		l.fields.Delete(key)
		fieldsLen++
		return true
	})

	if len(messages) == 0 || fieldsLen == 0 {
		return
	}

	entry := l.logger.WithFields(tempLoggerFields)
	entry.Tracef("%+v", messages[0].Message)
	l.fields.Store("stack", []message{})
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
	stackVal, _ := l.fields.Load("stack")
	stack := ensureStackType(stackVal)
	l.fields.Store("stack", append(stack, msg))
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

func newLogger(serviceName string) (logger *Logger) {
	formatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
		// PrettyPrint:     true,
		FieldMap: log.FieldMap{
			log.FieldKeyMsg: "log_message",
		},
	}

	newLogger := newLog(formatter, os.Stdout, log.TraceLevel, false)

	logger = new(Logger)
	logger.logger = newLogger
	logger.fields.Range(func(key interface{}, value interface{}) bool {
		logger.fields.Delete(key)
		return true
	})
	logger.fields.Store("stack", []message{})

	id := uuid.NewV1().String()
	logger.fields.Store(FieldLogID, id)
	logger.fields.Store(FieldServiceName, serviceName)
	logger.id = id
	return
}

func NewLogger(serviceName string) (logger *Logger) {
	logger = newLogger(serviceName)
	return
}
