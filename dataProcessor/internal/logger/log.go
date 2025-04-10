package logger

import (
	"context"
	"log"
	"os"
)

type ErrLogger struct {
	file   *os.File
	logger *log.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

func NewErrLogger(path string) (*ErrLogger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	res := &ErrLogger{
		file:   file,
		logger: log.New(file, "", log.Ldate|log.Ltime),
		ctx:    ctx,
		cancel: cancel,
	}

	return res, err
}

func (l *ErrLogger) LogError(mess string, err error) {
	select {
	case <-l.ctx.Done():
		return
	default:
		l.logger.Println(mess, " : ", err.Error())
	}

}

func (l *ErrLogger) Close() error {
	l.cancel()
	if err := l.file.Close(); err != nil {
		return err
	}
	return nil
}
