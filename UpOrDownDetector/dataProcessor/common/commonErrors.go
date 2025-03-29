package common

type EntityError byte

const (
	ProcessorError = EntityError(iota)
	ConsumerError
	SaverError
	NotificationError
)

type ErrorInfo struct {
	Entity EntityError
	Err    error
}
