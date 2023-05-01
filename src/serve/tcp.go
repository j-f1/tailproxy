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

	ip, err := net.ResolveIPAddr("ip", config.Target.Hostname())
	if err != nil {
		logger.Fatal("invalid target hostname: %v", err)
		return err
	}

	addr := &net.TCPAddr{
		IP:   ip.IP,
		Port: port,
	}

	listener := ts.ListenTailnet(port)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Fatal("error accepting connection: %v", err)
		}
		go handleTCP(conn, addr)
	}
}

func handleTCP(conn net.Conn, addr *net.TCPAddr) {
	defer conn.Close()
	proxied, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Err("error connecting to target: %v", err)
		return
	}
	defer proxied.Close()
	go io.Copy(proxied, conn)
	io.Copy(conn, proxied)
}
