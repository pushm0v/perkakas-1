package structs

import "encoding/json"

// type StatusError struct {
// 	Code       int
// 	Err        error
// 	ThirdParty bool
// }

// func (se StatusError) Error() string {
// 	return se.Err.Error()
// }

// func (se StatusError) Status() int {
// 	return se.Code
// }

// APIError defines error data format
type APIError struct {
	HTTPStatus int    `json:"http_status"`
	ThirdParty bool   `json:"-"`
	ID         string `json:"id"`
	EN         string `json:"en"`
}

func (e APIError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
