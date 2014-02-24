package base

// import (
// 	"base/socket"
// 	"log"
// )

// var (
// 	socketSrv *socket.Server
// )

// func main() {
// 	config := socket.NewConfig()
// 	config.CloseingTimeout = 10000
// 	config.Addr = "127.0.0.1:9000"
// 	config.CodecFactory = socket.NewDefaultCodecFactory()
// 	config.ConnectedHandler = connectedHandler
// 	config.DisconnectHandler = disconnectedHandler
// 	config.MessageHandler = messageHandler

// 	var err error
// 	socketSrv, err = socket.NewServer(config)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	socketSrv.Start()

// }

// func messageHandler(channel socket.IChannel, protoPack *socket.ProtoPack) {
// 	log.Println("Id:", protoPack.Id)
// 	log.Println("Iscompressed:", protoPack.Iscompressed)
// 	log.Println("Isencryted:", protoPack.Isencrypted)
// 	log.Println("PlatformId:", protoPack.PlatformId)
// 	log.Println("Body:", protoPack.Body)

// 	channel.Close()
// }

// func connectedHandler(channel socket.IChannel) {
// 	log.Println("客户端已连接。")
// }

// func disconnectedHandler(channel socket.IChannel) {
// 	log.Println("客户端已断开连接。")
// }
