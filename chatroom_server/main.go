package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const port = ":8080"

type Message struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type Client struct {
	name string
	conn net.Conn
}

var (
	messages = make(chan Message)
)

func main() {
	// 创建一个本地地址，作为聊天室的地址
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("监听端口失败：", err)
		return
	}
	defer listener.Close()
	fmt.Printf("服务器启动成功，监听端口为%v, 等待客户端连接...\n", port)

	clients := make(map[Client]bool)

	go broadcaster(clients)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("建立连接失败：", err)
			continue
		}
		go handleConn(conn, clients)
	}
}

// broadcaster 广播器，负责接收客户端消息并进行广播
func broadcaster(clients map[Client]bool) {
	for {
		msg := <-messages
		switch msg.Type {
		case "enter":
			msg.Text = msg.Name + "进入了聊天室"
		case "leave":
			msg.Text = msg.Name + "离开了聊天室"
		case "msg":
			msg.Text = msg.Name + "：" + msg.Text
		}
		// 广播消息给所有客户端
		for cli := range clients {
			cli.conn.Write([]byte(msg.Text + "\n"))
		}
	}
}

func handleConn(conn net.Conn, clients map[Client]bool) {
	var msg Message
	input := bufio.NewScanner(conn)
	for input.Scan() {
		err := json.Unmarshal(input.Bytes(), &msg)
		if err != nil {
			// 如果输入不是有效的JSON，记录错误并继续
			log.Printf("无法解析JSON消息：%v", err)
			continue
		}
		switch msg.Type {
		case "enter":
			clients[Client{msg.Name, conn}] = true
		case "leave":
			delete(clients, Client{msg.Name, conn})
		}
		messages <- msg
	}

	conn.Close()
}
