package socket

import (
	"encoding/binary"
	"io"
	//"log"
	"math"
	"sync"
)

/**
 * 编码解码接口
 * @author abram
 * @since 0.0.1
 */
type ICodec interface {
	/**
	 * 解码
	 * @param buf socket 接受到的二进制数据
	 * @author abram
	 * @return protoPack 数据包对象
	 * @return err 错误
	 */
	Decode() (protoPack *ProtoPack, err error)

	/**
	 * 编码
	 * @param protoPack 数据包对象
	 * @author abram
	 * @return buf 二进制数据流
	 * @return err
	 */
	Encode(protoPack ProtoPack) error

	WriteBool(value bool) error
	WriteByte(value byte) error
	WriteInt16(value int16) error
	WriteInt32(value int32) error
	WriteInt64(value int64) error
	WriteDouble(value float64) error
	WriteString(value string) error
	WriteBinary(value []byte) error
	Flush() error

	ReadBool() (bool, error)
	ReadByte() (byte, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadDouble() (float64, error)
	ReadString() (string, error)
	ReadBinary() ([]byte, error)

	Close() error
	FlushAndClose() error
}

/**
 * 默认的编码解码器
 * @author abram
 */
type DefaultCodec struct {
	lock      sync.RWMutex
	buffer    [8]byte
	transport ITransport //FramedTransport
}

func NewDefaultCodec(transport ITransport) ICodec {
	return &DefaultCodec{transport: transport}
}

/**
 * 默认的解码器
 * @author abram
 * @param data [] bytes
 * @return protoPack
 */
func (codec *DefaultCodec) Decode() (protoPack *ProtoPack, err error) {
	protoPack = NewProtoPack()
	v, err := codec.ReadByte()
	if err != nil {

	}
	protoPack.Isencrypted = v

	v, err = codec.ReadByte()
	if err != nil {
		return nil, err
	}
	protoPack.Iscompressed = v

	v, err = codec.ReadByte()
	if err != nil {
		return nil, err
	}
	protoPack.PlatformId = v

	var v16 int16
	v16, err = codec.ReadInt16()
	if err != nil {
		return nil, err
	}
	protoPack.Id = v16

	var bv []byte
	bv, err = codec.ReadBinary()
	if err != nil {
		return nil, err
	}
	protoPack.Body = bv
	return protoPack, nil
}

/**
 * 数据编码方法
 * @author abram
 * @param protoPack
 * @return []byte
 */
func (codec *DefaultCodec) Encode(protoPack ProtoPack) error {
	codec.lock.RLock()
	defer codec.lock.RUnlock()

	if err := codec.WriteByte(protoPack.Isencrypted); err != nil {
		return err
	}
	if err := codec.WriteByte(protoPack.Iscompressed); err != nil {
		return err
	}
	if err := codec.WriteByte(protoPack.PlatformId); err != nil {
		return err
	}
	if err := codec.WriteInt16(protoPack.Id); err != nil {
		return err
	}
	if err := codec.WriteBinary(protoPack.Body); err != nil {
		return err
	}
	if err := codec.Flush(); err != nil {
		return err
	}
	return nil
}

func (codec *DefaultCodec) ReadAll(buf []byte) error {
	_, err := io.ReadFull(codec.transport, buf)
	return err
}

func (codec *DefaultCodec) ReadByte() (byte, error) {
	buf := codec.buffer[0:1]
	err := codec.ReadAll(buf)
	return buf[0], err

}

func (codec *DefaultCodec) ReadBool() (bool, error) {
	b, err := codec.ReadByte()
	v := true
	if b != 1 {
		v = false
	}
	return v, err
}

func (codec *DefaultCodec) ReadInt16() (int16, error) {
	buf := codec.buffer[0:2]
	err := codec.ReadAll(buf)
	return int16(binary.BigEndian.Uint16(buf)), err
}

func (codec *DefaultCodec) ReadInt32() (int32, error) {
	buf := codec.buffer[0:4]
	err := codec.ReadAll(buf)
	return int32(binary.BigEndian.Uint32(buf)), err
}

func (codec *DefaultCodec) ReadInt64() (int64, error) {
	buf := codec.buffer[0:8]
	err := codec.ReadAll(buf)
	return int64(binary.BigEndian.Uint64(buf)), err
}

func (codec *DefaultCodec) ReadDouble() (float64, error) {
	buf := codec.buffer[0:8]
	err := codec.ReadAll(buf)
	return math.Float64frombits(binary.BigEndian.Uint64(buf)), err
}

func (codec *DefaultCodec) ReadString() (string, error) {
	size, err := codec.ReadInt32()
	if err != nil {
		return "", err
	}

	return codec.ReadStringBody(int(size))
}

func (codec *DefaultCodec) ReadStringBody(size int) (string, error) {
	if size < 0 {
		return "", nil
	}
	buf := make([]byte, size)
	_, err := io.ReadFull(codec.transport, buf)
	return string(buf), err
}

func (codec *DefaultCodec) ReadBinary() ([]byte, error) {
	size, err := codec.ReadInt32()
	if err != nil {
		return nil, err
	}
	isize := int(size)
	buf := make([]byte, isize)
	_, e := io.ReadFull(codec.transport, buf)
	return buf, e
}

func (codec *DefaultCodec) WriteByte(value byte) error {
	v := []byte{value}
	_, err := codec.transport.Write(v)
	return err
}

func (codec *DefaultCodec) WriteBool(value bool) error {
	if value {
		return codec.WriteByte(1)
	}
	return codec.WriteByte(0)
}

func (codec *DefaultCodec) WriteInt16(value int16) error {
	v := codec.buffer[0:2]
	binary.BigEndian.PutUint16(v, uint16(value))
	_, err := codec.transport.Write(v)
	return err
}

func (codec *DefaultCodec) WriteInt32(value int32) error {
	v := codec.buffer[0:4]
	binary.BigEndian.PutUint32(v, uint32(value))
	_, err := codec.transport.Write(v)
	return err
}

func (codec *DefaultCodec) WriteInt64(value int64) error {
	v := codec.buffer[0:8]
	binary.BigEndian.PutUint64(v, uint64(value))
	_, err := codec.transport.Write(v)
	return err
}

func (codec *DefaultCodec) WriteDouble(value float64) error {
	return codec.WriteInt64(int64(math.Float64bits(value)))
}

func (codec *DefaultCodec) WriteBinary(value []byte) error {
	err := codec.WriteInt32(int32(len(value)))
	if err != nil {
		return err
	}
	_, e := codec.transport.Write(value)
	return e
}

func (codec *DefaultCodec) WriteString(value string) error {
	return codec.WriteBinary([]byte(value))
}

func (codec *DefaultCodec) Flush() error {
	return codec.transport.Flush()
}

func (codec *DefaultCodec) Close() error {
	return codec.transport.Close()
}

func (codec *DefaultCodec) FlushAndClose() error {
	return codec.FlushAndClose()
}
