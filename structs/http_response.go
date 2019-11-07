package structs

// Meta defines meta format format for api format
type Meta struct {
	Version string `json:"version"`
	Status  string `json:"api_status"`
	APIEnv  string `json:"api_env"`
}

type Response struct {
	ResponseCode string       `json:"response_code"`
	ResponseDesc ResponseDesc `json:"response_desc"`
	Meta         Meta         `json:"meta"`
}

type SuccessResponse struct {
	Response
	Next *string     `json:"next,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// error Response
type ErrorResponse struct {
	Response
	HttpStatus int `json:"-"`
}

func (e *ErrorResponse) Error() string {
	return e.ResponseDesc.EN
}

// ResponseDesc defines details data response
type ResponseDesc struct {
	ID string `json:"id"`
	EN string `json:"en"`
}
