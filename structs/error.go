package structs

import (
	"net/http"
)

var ErrUnknown *ErrorResponse = &ErrorResponse{
	Response: Response{
		ResponseCode: "00001",
		ResponseDesc: ResponseDesc{
			ID: "Ups ada kesalahan, silahkan coba beberapa saat lagi",
			EN: "Unknown error",
		},
	},
	HttpStatus: http.StatusInternalServerError,
}

var ErrUnauthorized *ErrorResponse = &ErrorResponse{
	Response: Response{
		ResponseCode: "00001",
		ResponseDesc: ResponseDesc{
			ID: "Anda tidak diijinkan",
			EN: "You are not authorized",
		},
	},
	HttpStatus: http.StatusUnauthorized,
}
