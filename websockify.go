package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	cfg         Config
	err         error
	logger      *log.Logger
	host        string
	port        string
	proxyserver *ProxyServer
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	logger.Printf("new connection from: %s", r.RemoteAddr)
	wsConn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		logger.Fatal(err)
		return
	}

	dataType, data, err := wsConn.ReadMessage()
	if err != nil {
		_ = wsConn.WriteMessage(dataType, []byte("Fail: "+err.Error()))
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", string(data))
	if err != nil {
		errorMsg := "FAIL(net resolve tcp addr): " + err.Error()
		logger.Println(errorMsg)
		_ = wsConn.WriteMessage(websocket.CloseMessage, []byte(errorMsg))
		return
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		errorMsg := "FAIL(net dial tcp): " + err.Error()
		logger.Println(errorMsg)
		_ = wsConn.WriteMessage(websocket.CloseMessage, []byte(errorMsg))
		return
	}

	proxyserver = NewProxyServer(wsConn, tcpConn)
	go proxyserver.doProxy()
}

func main() {
	logger.Printf("listening on %s\n", cfg.Server.Addr())
	http.HandleFunc("/", handleConnection)
	http.ListenAndServe(cfg.Server.Addr(), nil)
}
