package base

// import (
// 	"base/socket"
// 	"log"
// )

// func main() {
// 	config := socket.NewConfig()
// 	config.Addr = "127.0.0.1:9000"
// 	config.CloseingTimeout = 0
// 	config.CodecFactory = socket.NewDefaultCodecFactory()
// 	config.MessageHandler = messageHandler
// 	config.ConnectedHandler = connectedHandler
// 	config.DisconnectHandler = disconnectedHandler

// 	client, err := socket.NewClient(config)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	err = client.Open()
// 	if err != nil {
// 		log.Println(err)
// 	}

// }

// func messageHandler(channel socket.IChannel, protoPack *socket.ProtoPack) {

// }

// func connectedHandler(channel socket.IChannel) {
// 	log.Println("服务端已连接。")

// 	protoPack := socket.NewProtoPack()
// 	protoPack.Id = 10
// 	protoPack.Iscompressed = 1
// 	protoPack.Isencrypted = 1
// 	protoPack.PlatformId = 12
// 	protoPack.Body = []byte{1, 1, 1, 1, 1}

// 	err := channel.Write(*protoPack)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	protoPack = socket.NewProtoPack()
// 	protoPack.Id = 11
// 	protoPack.Iscompressed = 1
// 	protoPack.Isencrypted = 1
// 	protoPack.PlatformId = 12
// 	protoPack.Body = []byte{2, 2, 2, 2, 2}

// 	err = channel.Write(*protoPack)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func disconnectedHandler(channel socket.IChannel) {
// 	log.Println("已从服务端断开连接。")
// }
