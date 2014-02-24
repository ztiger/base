package socket

import (
	"bytes"
	"encoding/binary"
	"io"
)

type FramedTransport struct {
	socket      ITransport // 实际类型为Socket
	writeBuffer *bytes.Buffer
	readBuffer  *bytes.Buffer
}

//socket 的世界类型为Socket
func NewFramedTransport(socket ITransport) *FramedTransport {
	writeBuf := make([]byte, 0, 1024)
	readBuf := make([]byte, 0, 1024)
	return &FramedTransport{socket: socket, writeBuffer: bytes.NewBuffer(writeBuf), readBuffer: bytes.NewBuffer(readBuf)}
}

func (transport *FramedTransport) Read(buf []byte) (int, error) {
	if transport.readBuffer.Len() > 0 {
		got, err := transport.readBuffer.Read(buf)
		if got > 0 {
			return got, err
		}
	}

	transport.readFrame()
	got, err := transport.readBuffer.Read(buf)
	return got, err
}

func (transport *FramedTransport) Write(buf []byte) (int, error) {
	n, err := transport.writeBuffer.Write(buf)
	return n, err
}

func (transport *FramedTransport) Flush() error {
	size := transport.writeBuffer.Len()
	buf := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(buf, uint32(size))

	if _, err := transport.socket.Write(buf); err != nil {
		return err
	}
	if size > 0 {
		if _, err := transport.writeBuffer.WriteTo(transport.socket); err != nil {
			return err
		}
	}
	err := transport.socket.Flush()
	return err
}

func (transport *FramedTransport) readFrame() (int, error) {
	buf := []byte{0, 0, 0, 0}
	if _, err := io.ReadFull(transport.socket, buf); err != nil {
		return 0, err
	}

	size := int(binary.BigEndian.Uint32(buf))
	if size == 0 {
		return 0, nil
	}

	buf2 := make([]byte, size)
	if n, err := io.ReadFull(transport.socket, buf2); err != nil {
		return n, err
	}
	transport.readBuffer = bytes.NewBuffer(buf2)
	return size, nil
}

func (transport *FramedTransport) Open() error {
	return transport.socket.Open()
}

func (transport *FramedTransport) IsOpen() bool {
	return transport.socket.IsOpen()
}

func (transport *FramedTransport) Peek() bool {
	return transport.socket.Peek()
}

func (transport *FramedTransport) FlushAndClose(flush bool) error {
	if err := transport.Flush(); err != nil {
		return err
	}
	return transport.socket.Close()
}

func (transport *FramedTransport) Close() error {
	return transport.socket.Close()
}
