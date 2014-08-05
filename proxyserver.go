package main

import (
	"net"

	"github.com/gorilla/websocket"
)

type ProxyServer struct {
	wsConn  *websocket.Conn
	tcpConn *net.TCPConn
}

func NewProxyServer(wsConn *websocket.Conn, tcpConn *net.TCPConn) *ProxyServer {
	proxyserver := ProxyServer{wsConn, tcpConn}
	return &proxyserver
}

func (proxyserver *ProxyServer) doProxy() {
	go proxyserver.wsToTcp()
	proxyserver.tcpToWs()
}

func (proxyserver *ProxyServer) tcpToWs() {
	buffer := make([]byte, 1024)

	for {
		n, err := proxyserver.tcpConn.Read(buffer)
		if err != nil {
			proxyserver.tcpConn.Close()
			break
		}

		err = proxyserver.wsConn.WriteMessage(websocket.BinaryMessage, buffer[0:n])
		if err != nil {
			logger.Println(err.Error())
		}
	}
}

func (proxyserver *ProxyServer) wsToTcp() {
	for {
		_, data, err := proxyserver.wsConn.ReadMessage()
		if err != nil {
			break
		}

		_, err = proxyserver.tcpConn.Write(data)
		if err != nil {
			logger.Println(err.Error())
			break
		}
	}
}
