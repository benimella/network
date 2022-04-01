package lib

import (
	"fmt"
	"log"
	"net"
	"time"
)

// address  ip:port
func Server(address string)  {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Received message %s => %s", conn.RemoteAddr(), conn.LocalAddr())

		go func(connection net.Conn) {
			for {
				buf := make([]byte, 1024)
				_, err := conn.Read(buf)
				if err != nil {
					log.Println("conn read err", err)
					break
				}
				_, err = conn.Write([]byte("ok"))
				if err != nil {
					log.Println("conn write err", err)
					break
				}
			}
		}(conn)
	}
}

// address  ip:port
func Client(address string)  {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		_, err = conn.Write([]byte("hello"))
		if err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second * 2)
	}
}
