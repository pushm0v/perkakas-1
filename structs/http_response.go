package structs

// Meta defines meta format format for api format
type Meta struct {
	Version     string `json:"version"`
	Status      string `json:"api_status"`
	APIEnv      string `json:"api_env"`
}

type Response struct {
	HTTPCode int         `json:"http_code"`
	Resp     interface{} `json:"response"`
	Meta     Meta        `json:"meta"`
}

type SuccessResponse struct {
	APICode int         `json:"api_code"`
	Next    *string     `json:"next,omitempty"`
	Data    interface{} `json:"data"`
}

// error Response
type ErrorResponse struct {
	APICode int       `json:"api_code"`
	Errors  ErrorData `json:"errors,omitempty"`
}

// ErrorData defines error data response
type ErrorData struct {
	Details DetailData `json:"details"`
}

// DetailData defines details data response
type DetailData struct {
	ID string `json:"id"`
	EN string `json:"en"`
}
