package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}
func (this *Server) hadle(conn net.Conn) {
	fmt.Println("链接建立成功")
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
	//for循环等待
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("服务器出错了：%s,继续下一个", err)
			continue
		}
		//go handle启动
		go this.hadle(conn)
	}
}
