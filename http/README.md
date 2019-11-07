# Package http

This package contains http handler related to help you write handler easier and have a standardized response.

## How to create custom handler
To create custom handler using this package, you need to declare the handler using this function signature:

```go
func(w http.ResponseWriter, r *http.Request) (data interface{}, pageToken *string, err error)
```

This handler requires you to return the `response data`, `page token` and `error` if any.
`response data` is the struct for the response data
`page token` is the optional hashed or encoded page token for pagination. 
We enforce this approach for pagination since it's the most clean way to do a pagination.'
This approach support infinite scroll and traditional pagination.
`error` is the error (if any)

Example:
```go
func HelloHandler(w http.ResponseWriter, r *http.Request) (data interface{}, pageToken *string, err error) {
    // Data we will return
	data = Person{
		FirstName: "Budi",
		LastName:  "Ariyanto",
	}

    // The token
	token := "o934ywjkhk67j78sd9af=="
    pageToken = &token
    
    // error, should be something like this
    // if err != nil {
    //    return err
    //}

	return
}

func main() {
    // M is meta
    // E is errorMap. Clients should supply the error map, so when error occured, the response writer
    // will render the correct error message, or it will render an unknown error.
	handlerCtx := phttp.HttpHandlerContext{
		M: structs.Meta{
			Version: "v1.2.3",
			Status:  "stable",
			APIEnv:  "prod-test",
		},
		E: errMap,
	}

    // newHandler is a function that will create function for create new custom handler with injected handler context
    newHandler := phttp.NewHttpHandler(handlerCtx)
    
    // helloHandler is the handler
	helloHandler := newHandler(HelloHandler)

	router := chi.NewRouter()
	router.Get("/hello", helloHandler.ServeHTTP)

	http.ListenAndServe(":5678", router)
}
```

From the example above, you can see you only care about the data, pageToken and error,
then this custom handler will construct the response itself.
This response are refer to [Kitabisa API response standardization](https://app.gitbook.com/@kitabisa-engineering/s/backend/standardization-1/api-response).
