package http

import (
	"errors"
	"net/http"

	"github.com/kitabisa/perkakas/v2/structs"
)

type HttpHandler struct {
	// C HttpHandlerContext
	H func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)
	CustomWriter
}

func NewHttpHandler(c HttpHandlerContext) func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)) HttpHandler {
	return func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, *string, error)) HttpHandler {
		return HttpHandler{H: handler, CustomWriter: CustomWriter{C: c}}
	}
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, pageToken, err := h.H(w, r)
	if err != nil {
		var apiError *structs.APIError
		if errors.As(err, &apiError) {
			h.WriteError(w, err)
		} else {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
		return
	}

	h.Write(w, data, pageToken)
}
