package structs

import (
	"net/http"
)

var ErrUnknown *ErrorResponse = &ErrorResponse{
	APICode:    "000001",
	HTTPStatus: http.StatusInternalServerError,
	Errors: ErrorData{
		Details: DetailData{
			ID: "Ups ada kesalahan, silahkan coba beberapa saat lagi",
			EN: "Unknown error",
		},
	},
}

var ErrUnauthorized *ErrorResponse = &ErrorResponse{
	APICode:    "000001",
	HTTPStatus: http.StatusUnauthorized,
	Errors: ErrorData{
		Details: DetailData{
			ID: "Anda tidak diijinkan",
			EN: "You are not authorized",
		},
	},
}
