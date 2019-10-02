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
		LogID:          "90744197-c2d9-46a0-b9cd-26e1fc2a4cbf", // expected to be uuid string, you can use other unique identifier
		Endpoint:       "http://example.com/query",
		Method:         "GET",
		RequestBody:    nil, // this should be filled when do POST request
		RequestHeader:  nil, // this should be filled when do POST request 
		ResponseBody:   `{"id": "gas976df97as6d", "name": "perkakas"}`, // get the value from resp.Body when you call http request, or response you give for your client
		ResponseHeader: `resp.Header`, // get the value from resp.Header
		ErrorMessage:   nil, // expected to be string or error type
	}
	logger.Set(fields)
	logger.Log(InfoLevel, "log message")
}
```

Example 2 - Create logger with builder (chaining):
```
	logger := NewLogger()
	 
	logger.
		SetLogID("90744197-c2d9-46a0-b9cd-26e1fc2a4cbf").
		SetEndpoint("http://example.com/query").
		SetMethod("GET").
		SetRequestBody(nil).
        SetRequestHeaders(nil).
		SetResponseBody(`{"id": "gas976df97as6d", "name": "perkakas"}`).
		SetResponseHeaders(resp.Header).
		SetErrorMessage(errors.New("Error in code 2123123"))
	logger.Set(fields)
	logger.Log(InfoLevel, "log message")

```