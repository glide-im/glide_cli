package apis

import "fmt"

type ApiError struct {
	Code    int
	Message string
}

func (a *ApiError) Error() string {
	return fmt.Sprintf("ApiError: %d-%s", a.Code, a.Message)
}
