package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.lisentChan()
	return user
}
func (this *User) lisentChan() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//发送msg到chan
	this.server.BroadCast(this, "上线")
}
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Addr)
	this.server.mapLock.Unlock()

	//发送msg到chan
	this.server.BroadCast(this, "下线")
}

func (this *User) DoMsg(msg string) {
	this.server.BroadCast(this, msg)
}
