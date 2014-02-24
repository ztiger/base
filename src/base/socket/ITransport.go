package socket

import (
	"io"
)

type Flusher interface {
	Flush() (err error)
}

type ITransport interface {
	io.ReadWriteCloser
	Flusher
	Open() error
	IsOpen() bool
	Peek() bool
}
