package socket

import (
	//"bytes"
	"errors"
	"log"
	"sync"
	"time"
)

var (
	tcp string = "tcp"
	ip4 string = "ip4"
)

const (
	Closing_timeout = 10000 // 关闭超时
	Listening_port  = 8888  //默认监听的端口
)

type Config struct {
	CloseingTimeout   time.Duration //关闭连接的超时时间
	Addr              string        //监听地址
	CodecFactory      ICodecFactory
	ConnectedHandler  func(channel IChannel)
	DisconnectHandler func(channel IChannel)
	MessageHandler    func(channel IChannel, protoPack *ProtoPack) //业务处理函数
}

/**
 * 生成一个socket 配置对象
 * @author abram
 */
func NewConfig() *Config {
	return &Config{}
}

/**
 * socket 服务对象
 * @author abram
 */
type Server struct {
	stopped          bool
	closingTimeout   time.Duration
	addr             string
	codecFactory     ICodecFactory
	mutex            sync.RWMutex
	serverSocket     *ServerSocket
	connectedHandler func(channel IChannel)
	disconnectHanler func(channel IChannel)
	messageHandler   func(channel IChannel, protoPack *ProtoPack)
}

/**
 * 生成SocketServer 对象
 * @param config
 * @return SocketServer
 */
func NewServer(config *Config) (*Server, error) {
	if config == nil {
		return nil, errors.New("config 不能为空。")
	}
	if config.ConnectedHandler == nil {
		return nil, errors.New("config.ConnectedHandler 不能为空。")
	}

	if config.DisconnectHandler == nil {
		return nil, errors.New("config.DisconnectHandler 不能为空。")
	}

	if config.MessageHandler == nil {
		return nil, errors.New("config.MessageHandler 不能为空。")
	}

	server := &Server{}

	server.closingTimeout = config.CloseingTimeout
	server.addr = config.Addr
	server.codecFactory = config.CodecFactory
	server.connectedHandler = config.ConnectedHandler
	server.messageHandler = config.MessageHandler
	server.disconnectHanler = config.DisconnectHandler

	if server.closingTimeout == 0 {
		server.closingTimeout = Closing_timeout
	}

	serverSocket, err := NewServerSocket(server.addr)
	if err != nil {
		return nil, err
	}
	server.serverSocket = serverSocket
	server.stopped = true
	return server, nil
}

/**
 * 启动服务器
 * @author abram
 */
func (server *Server) Start() error {
	if !server.stopped {
		return errors.New("服务已经启动。")
	}

	server.stopped = false
	err := server.serverSocket.Listen()
	if err != nil {
		return err
	}
	log.Println("开始监听...")
	for !server.stopped {
		client, err := server.serverSocket.Accept()
		if err != nil {
			log.Println("Accept err: ", err)
		}
		if client != nil {
			go func() {
				if err := server.connectionHandler(client); err != nil {
					log.Println("Error processing request:", err)
				}
			}()
		}
	}

	return nil
}

/**
 * 客户端接入管理
 * @author abram
 * @param client ITransport 实际类型是Socket
 */
func (server *Server) connectionHandler(client ITransport) error {
	transport := NewFramedTransport(client)
	codec := server.codecFactory.GetCodec(transport)
	channel := NewDefaultChannel(client, codec)

	defer func() {
		if server.disconnectHanler != nil {
			server.disconnectHanler(channel)
		}
		channel.Close()
		time.Sleep(500)
	}()

	server.connectedHandler(channel)
	for {
		protoPack, err := codec.Decode()
		if err != nil {
			break
		}
		go server.messageHandler(channel, protoPack)
	}

	return nil
}

/**
 * 关闭服务
 * @author abram
 */
func (server *Server) Stop() error {
	server.stopped = true
	server.serverSocket.Interrupt()
	return nil
}
