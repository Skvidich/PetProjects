package common

type EntityError byte

const (
	ProcessorError = EntityError(iota)
	ConsumerError
	ControllerError
	NotificationError
)

type ErrorInfo struct {
	Entity EntityError
	Err    error
}
