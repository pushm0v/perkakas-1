package httputil

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
)

const passwordPattern = `(\\{0,1}"password\\{0,1}"):\s*\\{0,1}"(.*?)\\{0,1}"`
var passRemover *regexp.Regexp

func init() {
	passRemover = regexp.MustCompile(passwordPattern)
}

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

// ExcludeSensitiveHeader exclude sensitive header. Currently, sensitive header only Authorization
func ExcludeSensitiveHeader(header http.Header) (h http.Header) {
	h = make(http.Header)
	for k, v := range header {
		h[k] = v
	}

	h.Del("Authorization")
	return
}

func ExcludeSensitiveRequestBody(body *string) {
	result := passRemover.ReplaceAllString(*body, "")
	*body = result
}