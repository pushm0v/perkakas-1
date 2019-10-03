# perkakas
perkakas/per·ka·kas/ - segala yang dapat dipakai sebagai alat (seperti untuk makan, bekerja di dapur, perang, ngoding)

Library for supporting common backend tasks. If you have function that considered common, please move it here.

How to develop on perkakas:
1. Any common function that **does not have any business logic** can be moved here.
2. You should group your function into a folder as a package.
3. If the folder not exist, just create it. If the folder exist, put your code there.

## Log
Example:
```
logger := NewLogger()

fields := Field{
	LogID:          suite.logID,
	Endpoint:       suite.url,
	Method:         "GET",
	RequestBody:    nil,
	RequestHeader:  nil,
	ResponseBody:   string(b),
	ResponseHeader: resp.Header,
	Message:        "Error",
	Level:          ErrorLevel,
}

logger.Set(fields).Print("testlog")	 
```

Example 2 - Create logger with builder (chaining):
```
logger := NewLogger()
	
logger.SetLoggerLevel(InfoLevel)
logger.
	SetLogID(suite.logID).
	SetEndpoint(suite.url).
	SetMethod("GET").
	SetRequestBody(nil).
	SetRequestHeaders(nil).
	SetResponseBody(string(b)).
	SetResponseHeaders(resp.Header).
	SetMessage(InfoLevel, errors.New("Error in code 2123123")).
	Print("testlog")
```