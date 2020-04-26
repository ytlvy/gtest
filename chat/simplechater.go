package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

//SimpleChater 简单聊天程序
type SimpleChater struct {
}

//Start 开启服务
func (s *SimpleChater) Start() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go s.broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handleConn(conn)
	}
}

type messageInfo struct {
	cliid int64
	msg   string
}

type client struct {
	id int64
	ch chan *messageInfo
}

var (
	enteringChan       = make(chan *client)
	leavingChan        = make(chan *client)
	messageChan        = make(chan *messageInfo)
	gClientid    int64 = 10000
)

func (s *SimpleChater) broadcast() {
	clients := make(map[*client]bool)
	for {
		select {
		case msg := <-messageChan:
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-enteringChan:
			clients[cli] = true

		case cli := <-leavingChan:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func (s *SimpleChater) makeClient() *client {
	ch := make(chan *messageInfo)
	gClientid++
	return &client{gClientid, ch}
}

func (s *SimpleChater) handleConn(conn net.Conn) {
	cli := s.makeClient()
	go s.clientWriter(conn, cli)

	who := conn.RemoteAddr().String()
	cli.ch <- &messageInfo{cli.id, "You are " + who}
	messageChan <- &messageInfo{cli.id, who + " has arrived"}
	enteringChan <- cli

	input := bufio.NewScanner(conn)
	for true {
		// fmt.Printf("input a msg:")
		input.Scan()
		text := input.Text()
		if err := input.Err(); err != nil {
			fmt.Println("Erro reading from input: ", err)
		}

		if text == "" {
			break
		}

		messageChan <- &messageInfo{cli.id, text}

	}

	leavingChan <- cli
	messageChan <- &messageInfo{cli.id, who + "has left"}
	conn.Close()
}

func (s *SimpleChater) clientWriter(conn net.Conn, cli *client) {

	for msgInfo := range cli.ch {
		who := conn.RemoteAddr().String()
		if cli.id == msgInfo.cliid {
			who = "you"
		}
		fmt.Fprintln(conn, who+" said: "+msgInfo.msg)
	}
}
