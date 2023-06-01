package serverModule

import (
	"awesomeProject/userModule"
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户map
	OnlineMap map[string]*userModule.User
	MapLock   sync.RWMutex
	Message   chan string
}

func (this *Server) Handler(conn net.Conn) {
	user := userModule.NewUser(conn)

	this.MapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.MapLock.Unlock()

	//广播上线
	this.Online(user)

	go func() {
		byt := make([]byte, 4096)
		for {
			n, err := conn.Read(byt)
			if n == 0 {
				this.Offline(user)
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("连接读取错误")
				return
			}
			msg := string(byt[:n-1])
			fmt.Println(msg)
			this.BroadCastMsg(user, msg)
		}
	}()

	//阻塞
	select {}
}

//广播功能
func (this *Server) BroadCastMsg(user *userModule.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

//监听Message
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		// 发给在线用户
		this.MapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.MapLock.Unlock()
	}
}

func (this *Server) Online(user *userModule.User) {
	this.MapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.MapLock.Unlock()
	//广播
	this.BroadCastMsg(user, "上线了")
}

func (this *Server) Offline(user *userModule.User) {
	this.MapLock.Lock()
	delete(this.OnlineMap, user.Name)
	this.MapLock.Unlock()
	//广播
	this.BroadCastMsg(user, "下线了")
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
