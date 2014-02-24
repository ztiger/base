package socket

import ()

//解码工厂接口
type ICodecFactory interface {
	GetCodec(transport ITransport) ICodec
}

type DefaultCodecFactory struct {
}

//获取默认的解码工厂类
func NewDefaultCodecFactory() ICodecFactory {
	return &DefaultCodecFactory{}
}

//获取默认的解码器
func (factory *DefaultCodecFactory) GetCodec(transport ITransport) ICodec {
	return NewDefaultCodec(transport)
}
