package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/kitabisa/perkakas/v2/structs"
)

type HttpHandlerContext struct {
	M structs.Meta
	E map[error]*structs.ErrorResponse
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

	successResp.APICode = "000000"
	successResp.Next = nextPage

	apiResponse := &structs.Response{
		Resp: successResp,
		Meta: c.C.M,
	}

	res, err := json.Marshal(apiResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (c *CustomWriter) WriteError(w http.ResponseWriter, err error) {
	if len(c.C.E) > 0 {
		errorResponse := LookupError(c.C.E, err)
		c.writeResponse(w, errorResponse)
	} else {
		var errorResponse *structs.ErrorResponse
		if errors.As(err, &errorResponse) {
			c.writeResponse(w, errorResponse)
		} else {
			c.writeResponse(w, structs.ErrUnknown)
		}
	}
}

func (c *CustomWriter) writeResponse(w http.ResponseWriter, errorResponse *structs.ErrorResponse) {
	apiResponse := &structs.Response{
		Resp: errorResponse,
		Meta: c.C.M,
	}

	res, err := json.Marshal(apiResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorResponse.HTTPStatus)
	w.Write(res)
}

// GetErrorMessage will get error message based on error type
func LookupError(lookup map[error]*structs.ErrorResponse, err error) (res *structs.ErrorResponse) {
	if msg, ok := lookup[err]; ok {
		res = msg
		return
	}

	return
}
