package rhttp

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

func NewServer(port int) error {
	listener, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(port)))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(conn.LocalAddr(), " - accepted connection")

		go serveConn(conn)
	}
}

func serveConn(conn net.Conn) {
	defer conn.Close()

	listener, err := net.Listen("tcp", "0.0.0.0:")
	if err != nil {
		log.Println(conn.LocalAddr(), " - failed to open port, ", err)
		return
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Println(conn.LocalAddr(), " - failed to close listener, ", err)
		}
	}()

	_, err = conn.Write([]byte(listener.Addr().String()))
	if err != nil {
		log.Println(conn.LocalAddr(), " - ", err)
		return
	}

	remote, err := listener.Accept()
	if err != nil {
		log.Println(conn.LocalAddr(), " - failed to accept connection, ", err)
		return
	}
	defer remote.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer conn.Close()
		defer wg.Done()
		if _, err := io.Copy(conn, remote); err != nil {
			log.Println(conn.LocalAddr(), " - failed copy outgoing traffic, ", err)
			return
		}
	}()

	go func() {
		defer remote.Close()
		defer wg.Done()
		if _, err := io.Copy(remote, conn); err != nil {
			log.Println(conn.LocalAddr(), " - failed copy incoming traffic, ", err)
			return
		}
	}()

	wg.Wait()
}
