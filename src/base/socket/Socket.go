package socket

import (
	"errors"
	"net"
	"time"
)

//socket 结构
type Socket struct {
	conn    net.Conn
	addr    net.Addr
	timeout time.Duration
}

//创建一个客户端的一个socket 连接
// hostPort 格式为host:port
func NewSocket(hostPort string) (*Socket, error) {
	return NewSocketTimeout(hostPort, 0)
}

//根据hostPort创建一个会超时的socket连接
//hostPort 格式 host:port
//timeout 超时时间
func NewSocketTimeout(hostPort string, timeout time.Duration) (*Socket, error) {
	addr, err := net.ResolveTCPAddr("tcp", hostPort)
	if err != nil {
		return nil, err
	}
	return NewSocketFromAddrTimeout(addr, timeout)
}

//根据net.Addr 常见一个socket 连接
func NewSocketFromAddrTimeout(addr net.Addr, timeout time.Duration) (*Socket, error) {
	return &Socket{addr: addr, timeout: timeout}, nil
}

//根据net.Conn 创建一个socket 对象，此方法用于服务端接受到Conn后使用
func NewSocketFromConnTimeout(conn net.Conn, timeout time.Duration) (*Socket, error) {
	return &Socket{conn: conn, timeout: timeout}, nil
}

//设置超时时间
func (socket *Socket) SetTimeout(timeout time.Duration) error {
	socket.timeout = timeout
	return nil
}

//使用超时机制
func (socket *Socket) pushDeadline(read, write bool) {
	var t time.Time
	if socket.timeout > 0 {
		t = time.Now().Add(time.Duration(socket.timeout))
	}
	if read && write {
		socket.conn.SetDeadline(t)
	} else if read {
		socket.conn.SetReadDeadline(t)
	} else if write {
		socket.conn.SetWriteDeadline(t)
	}
}

// 判断客户端socket是否打开
func (socket *Socket) IsOpen() bool {
	if socket.conn != nil {
		return true
	}
	return false
}

// 打开客户端的socket
func (socket *Socket) Open() error {
	if socket.IsOpen() {
		return nil
	}

	if socket.addr == nil || len(socket.addr.String()) == 0 {
		return errors.New("addr 为空。")
	}
	if len(socket.addr.Network()) == 0 {
		return errors.New("网络不好。")
	}

	var err error
	if socket.conn, err = net.DialTimeout(socket.addr.Network(), socket.addr.String(), socket.timeout); err != nil {
		return err
	}
	return nil
}

//获取net.Conn
func (socket *Socket) Conn() net.Conn {
	return socket.conn
}

//关闭连接
func (socket *Socket) Close() error {
	if socket.conn == nil {
		return nil
	}

	if err := socket.conn.Close(); err != nil {
		return err
	}
	socket.conn = nil
	return nil
}

//读取数据
func (socket *Socket) Read(buf []byte) (int, error) {
	if !socket.IsOpen() {
		return 0, errors.New("Socket 连接已关闭。")
	}

	socket.pushDeadline(true, false)
	n, err := socket.conn.Read(buf)
	return n, err
}

//写数据
func (socket *Socket) Write(buf []byte) (int, error) {
	if !socket.IsOpen() {
		return 0, errors.New("Socket 连接已关闭。")
	}

	socket.pushDeadline(false, true)
	n, err := socket.conn.Write(buf)
	return n, err
}

//中断连接
func (socket *Socket) Interrupt() error {
	if !socket.IsOpen() {
		return nil
	}
	return socket.conn.Close()
}

func (socket *Socket) Flush() error {
	return nil
}

func (socket *Socket) Peek() bool {
	return socket.IsOpen()
}
