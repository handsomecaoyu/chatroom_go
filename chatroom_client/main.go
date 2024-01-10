package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Text string `json:"text"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// 启动一个 goroutine 来监听来自服务器的消息。
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	name := getName()
	// 发送用户名信息
	msg := Message{name, "enter", ""} // 消息类型为 name
	sendMsg(conn, msg)

	// 读取标准输入并将消息发送到服务器。
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		if strings.TrimSpace(text) == "/quit" {
			break
		}
		msg := Message{name, "msg", text}
		fmt.Print("\033[1A")
		fmt.Print("\r\033[K")
		sendMsg(conn, msg)

	}
}

func getName() string {
	fmt.Print("给自己取一个昵称：")
	var name string
	_, err := fmt.Scanln(&name)
	if err != nil {
		fmt.Printf("读取名字时出错: %v\n", err)
		return "匿名"
	}
	fmt.Printf("你好, %s!\n", name)
	return name
}

func sendMsg(conn net.Conn, msg Message) {
	// 发送包含名字和消息的json信息
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonMsg = append(jsonMsg, '\n') // 添加一个换行符，以便服务器可以使用 ReadString 读取整行。
	_, err = conn.Write(jsonMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
