package baseerror

/*
通用错误error
*/

type (
	BaseError struct {
		message string
	}
)

func NewBaseError(message string) *BaseError {
	return &BaseError{message: message}
}

func (e *BaseError) Error() string {

	return e.message
}
