package logger

type Logger interface {
	LogError(mess string, err error)
}
