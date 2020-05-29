package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"tcp-heartbeat/message"
	"time"
)

var (
	port         int
	scanInterval int
	ttl          int
)

type client struct {
	name          string
	lastUpdatedAt time.Time
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.IntVar(&port, "p", 9090, "server port")
	flag.IntVar(&scanInterval, "s", 10, "scan interval seconds")
	flag.IntVar(&ttl, "S", 10, "conn expire seconds")
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("server listening on port: %d\n", port)

	//sync.Map没有len方法 https://github.com/golang/go/issues/20680
	conns := new(sync.Map)

	go scan(conns)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn, conns)
	}
}

func handleConn(conn net.Conn, conns *sync.Map) {
	bs := make([]byte, 1024)
	for {
		//read message from client
		n, err := conn.Read(bs)

		if err == nil {
			msg := message.Message{}
			json.Unmarshal(bs[:n], &msg)

			//set heartbeat time
			if msg.IsHeartBeat() {
				log.Println(msg.Owner + " says: " + msg.Content)

				conns.Store(conn, client{
					name:          msg.Owner,
					lastUpdatedAt: time.Now(),
				})
			}
		}
	}
}

//kickoff expired connections
func scan(conns *sync.Map) {
	for {
		log.Println("scanning...")

		conns.Range(func(k, v interface{}) bool {
			client := v.(client)
			if time.Now().Sub(client.lastUpdatedAt).Seconds() > float64(ttl) {
				k.(net.Conn).Close()
				conns.Delete(k)
				log.Printf("client (%s) was kicked off\n", client.name)
			}
			return true
		})

		time.Sleep(time.Duration(scanInterval) * time.Second)
	}
}
