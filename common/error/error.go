package error

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRecordNotFound      = NewError(codes.NotFound, "record not found")
	ErrInternalServerError = NewError(codes.Internal, "internal server error")
	ErrBadRequest          = NewError(codes.InvalidArgument, "bad request")
)

type Error struct {
	Code    codes.Code `json:"code"`
	Message string     `json:"message"`
}

func NewError(code codes.Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (err *Error) Error() error {
	return fmt.Errorf("%d:%s", err.Code, err.Message)
}

func ParseError(err error) *Error {
	var statusCode int32
	if st, ok := status.FromError(err); ok {
		statusCode = st.Proto().Code
	}

	var strToCode = map[int]codes.Code{
		0:  codes.OK,
		1:  codes.Canceled,
		2:  codes.Unknown,
		3:  codes.InvalidArgument,
		4:  codes.DeadlineExceeded,
		5:  codes.NotFound,
		6:  codes.AlreadyExists,
		7:  codes.PermissionDenied,
		8:  codes.ResourceExhausted,
		9:  codes.FailedPrecondition,
		10: codes.Aborted,
		11: codes.OutOfRange,
		12: codes.Unimplemented,
		13: codes.Internal,
		14: codes.Unavailable,
		15: codes.DataLoss,
		16: codes.Unauthenticated,
	}

	if err == nil {
		return nil
	}

	parts := strings.Split(err.Error(), "desc = ")

	errlen := len(parts)

	outputmsg := ""

	if errlen > 1 {
		outputmsg = parts[errlen-1]
	} else {
		outputmsg = "No Error Desc"
	}

	return NewError(strToCode[int(statusCode)], outputmsg)
}
