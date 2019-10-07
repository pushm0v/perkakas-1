# Log
This logger will help you logging on a request - response cycle and show the stack trace of the error when you add log message.
Logging level not supported. Level in this logger only for showing the severity level when add log messages.

`NewLogger()` will create a logger for you.

`SetRequest(req interface{})` is function expecting *http.Request, then will extract required information from the request struct to fill the log.
Later will support graphql and grpc also.
Input can also a primitive type, and should be string. It will treat the input as the request body in string.

`SetResponse(resp interface{}, body interface{})` will extract the information for the log from resp, expecting *http.Response.
Body will fill the response body in the log message.

`AddMessage(lv Level, msg interface{})` add message to the logger along with severity level.

`Print()` print the log message, then flush. When no message added, it will not print anything.

Example:
```
logger := NewLogger()
	
logger.SetRequest(req) // req expected as *http.Request
logger.SetResponse(resp, string(b)) // resp expectes as *http.Response
logger.AddMessage(TraceLevel, "This is trace")
logger.AddMessage(DebugLevel, "This is debug")
logger.AddMessage(InfoLevel, "This is info")
logger.AddMessage(WarnLevel, "This is warning")
logger.AddMessage(ErrorLevel, errors.New("This is error"))
logger.AddMessage(FatalLevel, "This is fatal")
logger.AddMessage(PanicLevel, "This is panic")
logger.Print()
```