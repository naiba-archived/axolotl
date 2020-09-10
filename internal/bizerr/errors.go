package bizerr

import "fmt"

var (
	UnAuthorizedError = BizError{Code: 10001, Msg: "You have to login to continue"}
	DatabaseError     = BizError{Code: 10100, Msg: "Database error"}
	UnknownError      = BizError{Code: 55555}
)

type BizError struct {
	Code uint
	Msg  string
}

func (e BizError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Msg)
}
