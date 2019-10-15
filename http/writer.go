package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/kitabisa/perkakas/v2/structs"
)

type HttpHandlerContext struct {
	M           structs.Meta
	SuccessCode int
	E           map[error]structs.ErrorResponse
}

type CustomWriter struct {
	C HttpHandlerContext
}

func (c *CustomWriter) Write(w http.ResponseWriter, data interface{}, nextPage *string) {
	var successResp structs.SuccessResponse
	voData := reflect.ValueOf(data)
	arrayData := []interface{}{}

	if voData.Kind() != reflect.Slice {
		if voData.IsValid() {
			arrayData = []interface{}{data}
		}
		successResp.Data = arrayData
	} else {
		if voData.Len() != 0 {
			successResp.Data = data
		} else {
			successResp.Data = arrayData
		}
	}

	successResp.APICode = c.C.SuccessCode
	successResp.Next = nextPage

	apiResponse := &structs.Response{
		HTTPCode: http.StatusOK,
		Resp:     successResp,
		Meta:     c.C.M,
	}

	res, err := json.Marshal(apiResponse)
	if err != nil {
		// handle error
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (c *CustomWriter) WriteError(w http.ResponseWriter, err error) {
	var apiError *structs.APIError
	if errors.As(err, &apiError) {
		errorResponse := LookupError(c.C.E, err)
		apiResponse := &structs.Response{
			HTTPCode: apiError.HTTPStatus,
			Resp:     errorResponse,
			Meta:     c.C.M,
		}

		res, err := json.Marshal(apiResponse)
		if err != nil {
			// handle error
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiError.HTTPStatus)
		w.Write(res)
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}
}

// GetErrorMessage will get error message based on error type
func LookupError(lookup map[error]structs.ErrorResponse, err error) (res structs.ErrorResponse) {
	if msg, ok := lookup[err]; ok {
		res = msg
		return
	}

	return
}
