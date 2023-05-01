package serve

import (
	"io"
	"net"
	"strconv"
	"tailproxy/src/config"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func ServeTCP() error {
	port, err := strconv.Atoi(config.Target.Port())
	if err != nil {
		logger.Fatal("invalid target port: %v", err)
		return err
	}

	listener := ts.ListenTailnet(port)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Err("error accepting connection: %v", err)
		} else {
			go handleTCP(conn)
		}
	}
}

func handleTCP(conn net.Conn) {
	defer conn.Close()
	destConn, err := net.Dial("tcp", config.Target.Host)
	if err != nil {
		logger.Err("error connecting to target: %v", err)
		return
	}
	defer destConn.Close()
	go io.Copy(destConn, conn)
	io.Copy(conn, destConn)
	select {}
}
