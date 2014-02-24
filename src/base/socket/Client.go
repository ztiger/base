package socket

import (
	"errors"
	//"fmt"
	//"log"
	"sync"
	"time"
)

type Client struct {
	stopped           bool
	closingTimeout    time.Duration
	addr              string
	codecFactory      ICodecFactory
	mutex             sync.RWMutex
	socket            *Socket
	connectedHandler  func(channel IChannel)                       //连接建立事件
	disconnectHandler func(channel IChannel)                       //连接断开事件
	messageHandler    func(channel IChannel, protoPack *ProtoPack) //消息处理逻辑
}

// 生成一个客户端对象
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, errors.New("Socket config 不能为空。")
	}

	if config.Addr == "" || len(config.Addr) == 0 {
		return nil, errors.New("config.addr 不能为空.")
	}

	if config.CodecFactory == nil {
		return nil, errors.New("config.DecodeFactory 不能为空。")
	}

	if config.ConnectedHandler == nil {
		return nil, errors.New("config.ConnectedHandler 不能为空。")
	}

	if config.MessageHandler == nil {
		return nil, errors.New("config.messageHandler 不能为空。")
	}
	if config.DisconnectHandler == nil {
		return nil, errors.New("config.disconnectHandler 不能为空。")
	}

	client := &Client{}
	client.addr = config.Addr
	client.closingTimeout = config.CloseingTimeout
	client.codecFactory = config.CodecFactory
	client.connectedHandler = config.ConnectedHandler
	client.messageHandler = config.MessageHandler
	client.disconnectHandler = config.DisconnectHandler

	client.stopped = true
	return client, nil
}

//
func (client *Client) Open() error {
	if client.stopped == false {
		return errors.New("Client 已经打开。")
	}

	client.stopped = false
	socket, err := NewSocketTimeout(client.addr, client.closingTimeout)
	if err != nil {
		return err
	}
	client.socket = socket
	if err := client.socket.Open(); err != nil {
		return err
	}
	client.connectionHandler()
	return nil
}

//判断client是否已经打开
func (client *Client) IsOpen() bool {
	if client.socket == nil {
		return false
	}
	return client.socket.IsOpen()
}

//关闭连接
func (client *Client) Close() error {
	if client.socket == nil {
		return nil
	}

	return client.socket.Close()
}

//处理连接
func (client *Client) connectionHandler() error {
	transport := NewFramedTransport(client.socket)
	codec := client.codecFactory.GetCodec(transport)
	channel := NewDefaultChannel(client.socket, codec)

	defer func() {
		if client.disconnectHandler != nil {
			go client.disconnectHandler(channel)

		}
		channel.Close()

		time.Sleep(1000)
	}()
	client.connectedHandler(channel)
	var protoPack *ProtoPack
	var err error
	for {
		protoPack, err = codec.Decode()
		if err != nil {
			break
		}
		go client.messageHandler(channel, protoPack)
	}

	return err
}
