package custom_err

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"runtime"
	"strings"
)

var (
	ErrUnexpected                = errors.New("unexpected")
	ErrUnprocessableEntity       = errors.New("unprocessable entity")
	ErrSecondNumberLessThanFirst = errors.New("second number should be larger or equal")
	ErrMaxIntFibonacci           = errors.New("second number should be less than: 5000, use method for large numbers")
)

type CustomErrorWithCode struct {
	RestCode int
	GrpcCode codes.Code
	Msg      error
}

func New(originalError error, forUserError error, httpCode int, grpcCode codes.Code) *CustomErrorWithCode {
	return &CustomErrorWithCode{
		Msg:      fmt.Errorf("original error msg: %v, for user error message: %w", originalError, forUserError),
		RestCode: httpCode,
		GrpcCode: grpcCode,
	}
}

func (e *CustomErrorWithCode) Error() string {
	return e.Msg.Error()
}

func (e *CustomErrorWithCode) ErrorForUser() string {
	return UnwrapRecursive(e.Msg).Error()
}

func GetInfo() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	fn := runtime.FuncForPC(pc).Name()
	funcName := fn[strings.LastIndex(fn, ".")+1:]

	return fmt.Sprintf("%v: %%w", funcName)
}

func UnwrapRecursive(err error) error {
	unwrappedError := errors.Unwrap(err)
	if unwrappedError != nil {
		return UnwrapRecursive(unwrappedError)
	}
	return err
}
