package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
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
