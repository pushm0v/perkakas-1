package httputil

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func ReadRequestBody(req *http.Request) (bodyString string) {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}

	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	bodyString = string(bodyBytes)
	return
}
