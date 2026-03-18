package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	c         chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		c:         make(chan string),
	}
}
func (this *Server) LisentMsg() {
	for {
		msg := <-this.c
		fmt.Println("广播消息:", msg)
		this.mapLock.RLock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.RUnlock()
	}
}
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + msg
	this.c <- sendMsg
}
func (this *Server) hadle(conn net.Conn) {
	// fmt.Println("链接建立成功")
	user := NewUser(conn, this)
	user.Online()
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err", err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMsg(msg)
		}
	}()
	select {} //让当前协程 永久阻塞、不退出！让他发送完自己的上线也能继续听别人的上线
}
func (this *Server) Start() {
	//创建tcp连接器
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("创建tcp连接器出现异常:", err)
		return
	}
	//defer关闭
	defer listen.Close()
	go this.LisentMsg()
	//for循环等待
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		//go handle启动
		go this.hadle(conn)
	}
}
