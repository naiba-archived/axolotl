package bizerr

import "fmt"

var (
	UnAuthorizedError = BizError{Code: 10001}
	UnknownError      = BizError{Code: 55555}
)

type BizError struct {
	Code uint
}

func (e BizError) Error() string {
	return fmt.Sprintf("Error: %d", e.Code)
}
