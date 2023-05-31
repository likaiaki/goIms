package server_tool

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功")
}

func NewServer(ip string, port int) *Server {
	sever := &Server{
		Ip:   ip,
		Port: port,
	}
	return sever
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept error", err)
			continue
		}
		go this.Handler(conn)

	}
}
