package serverModule

import (
	"awesomeProject/userModule"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户map
	OnlineMap map[string]*userModule.User
	mapLock   sync.RWMutex
	Message   chan string
}

func (this *Server) Handler(conn net.Conn) {
	user := userModule.NewUser(conn)

	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//广播上线
	this.Br0adCastMsg(user, "上线了")

	//阻塞
	select {}
}

//广播功能
func (this *Server) Br0adCastMsg(user *userModule.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

//监听Message
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		// 发给在线用户
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func NewServer(ip string, port int) *Server {
	sever := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*userModule.User),
		Message:   make(chan string),
	}
	return sever
}

//启动服务
func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}
	defer listener.Close()
	//监听Message
	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept error", err)
			continue
		}
		go this.Handler(conn)

	}
}
