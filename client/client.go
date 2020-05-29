package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"tcp-heartbeat/message"
	"time"
)

var (
	host         string
	port         int
	name         string
	sendInterval int
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&host, "h", "127.0.0.1", "server host")
	flag.IntVar(&port, "p", 9090, "server port")
	flag.StringVar(&name, "n", randString(), "client name")
	flag.IntVar(&sendInterval, "s", 1, "heartbeat send interval")
}

func main() {
	flag.Parse()

	addr := host + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("client(%s) connected to %s\n", name, addr)

	defer conn.Close()

	go sendHeartbeat(conn)

	//wait forever
	ch := make(chan bool)
	ch <- true
}

func sendHeartbeat(conn net.Conn) {
	for {
		//send message to server
		msg := message.Message{
			MessageType: message.Heartbeat,
			Content:     "i am still alive",
			Owner:       name,
		}
		b, _ := json.Marshal(msg)
		conn.Write(b)
		
		log.Println("heartbeat was sent")
		
		time.Sleep(time.Duration(sendInterval) * time.Second)
	}
}

func randString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
