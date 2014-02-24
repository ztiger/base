package socket

import (
	"errors"
	"net"
	"time"
)

type ServerSocket struct {
	listener      net.Listener
	addr          net.Addr
	clientTimeout time.Duration
	interrupted   bool
}

func NewServerSocket(listenAddr string) (*ServerSocket, error) {
	return NewServerSocketTimeout(listenAddr, 0)
}

func NewServerSocketTimeout(listenAddr string, clientTimeout time.Duration) (*ServerSocket, error) {
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	return &ServerSocket{addr: addr, clientTimeout: clientTimeout}, nil
}

//判断是否已经在监听了
func (serverSocket *ServerSocket) IsListening() bool {
	if serverSocket.listener == nil {
		return false
	}
	return true
}

//开始监听
func (serverSocket *ServerSocket) Listen() error {
	if serverSocket.IsListening() {
		return errors.New("服务已经在监听了。")
	}
	l, err := net.Listen(serverSocket.addr.Network(), serverSocket.addr.String())
	if err != nil {
		return err
	}
	serverSocket.listener = l
	return nil
}

//接受客户端的请求
func (serverSocket *ServerSocket) Accept() (ITransport, error) {
	if serverSocket.interrupted {
		return nil, errors.New("Interrupted.")
	}
	if serverSocket.listener == nil {
		return nil, errors.New("Socket服务没打开。")
	}
	conn, err := serverSocket.listener.Accept()
	if err != nil {
		return nil, err
	}
	return NewSocketFromConnTimeout(conn, serverSocket.clientTimeout)
}

//获取监听地址
func (serverSocket *ServerSocket) Addr() net.Addr {
	return serverSocket.addr
}

//关闭服务
func (serverSocket *ServerSocket) Close() error {
	defer func() {
		serverSocket.listener = nil
	}()
	if serverSocket.IsListening() {
		serverSocket.listener.Close()
	}
	return nil
}

//中断服务
func (serverSocket *ServerSocket) Interrupt() error {
	serverSocket.interrupted = true
	return nil
}
