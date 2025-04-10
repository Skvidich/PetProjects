package errors

type ErrorType byte

const (
	TypeProcessor = ErrorType(iota)
	TypeConsumer
	TypeSaver
	TypeNotificator
)

type TaggedError struct {
	Type ErrorType
	Err  error
}
