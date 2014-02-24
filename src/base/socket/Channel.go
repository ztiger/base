package socket

import (
	"errors"
	"log"
)

type IChannel interface {
	Write(data interface{}) error
	SetAttribute(key string, val interface{})
	GetAttribute(key string) (interface{}, bool)
	Close() error
	IsOpen() bool
}

type DefaultChannel struct {
	codec      ICodec
	socket     ITransport
	attributes map[string]interface{}
}

func NewDefaultChannel(socket ITransport, codec ICodec) IChannel {
	return &DefaultChannel{socket: socket, codec: codec, attributes: make(map[string]interface{})}
}

func (channel *DefaultChannel) Write(data interface{}) error {
	log.Println("")
	if v, ok := data.(ProtoPack); ok {
		err := channel.codec.Encode(v)
		if err != nil {
			return errors.New("发送数据失败。")
		}
		return nil
	}

	return errors.New("错误的数据。")
}

// 关闭连接
func (channel *DefaultChannel) Close() error {
	return channel.codec.Close()
}

//把缓存中的数据写出去，然后再关闭连接
func (channel *DefaultChannel) FlushAndClose() error {
	return channel.codec.FlushAndClose()
}

func (channel *DefaultChannel) IsOpen() bool {
	return channel.socket.IsOpen()
}

func (channel *DefaultChannel) SetAttribute(key string, val interface{}) {
	channel.attributes[key] = val
}

func (channel *DefaultChannel) GetAttribute(key string) (interface{}, bool) {
	v, ok := channel.attributes[key]
	return v, ok
}
