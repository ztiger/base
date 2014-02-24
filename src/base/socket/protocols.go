package socket

/**
 * 消息体结构
 * @author abram
 */
type ProtoPack struct {
	Id           int16  //消息id
	Iscompressed byte   // 是否压缩 0-未压缩 1-压缩
	Isencrypted  byte   // 是否加密 0-未加密 1-加密
	PlatformId   byte   // 平台号
	Body         []byte // 消息体
}

/**
 * 生成一个ProtoPack 实例
 * @author abram
 * return ProtoPack
 */
func NewProtoPack() *ProtoPack {
	return &ProtoPack{}
}
