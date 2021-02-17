package server

import "net"

type connection struct {
	con net.Conn
	err error
}
