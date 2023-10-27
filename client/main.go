package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	c := new(Client) // 新客户端
	err := c.init()
	ErrorHandler(err)
	defer c.close() // 关闭客户端
	go cRun(c) // 运行客户端连接与接受程序
	runSend(c) // 客户端发送消息
}

/*
用于持续接受信息
*/
func receive(c *Client) {
	for {
		str, add, err := c.receiveInfo() // 接受信息
		if err != nil { //处理报错
			if strings.Contains(err.Error(),"not connected") {
				fmt.Println("Remote connected was closed. Or not connected to remote host.")
				break
			}
		}
		if err != nil && strings.Contains(err.Error(),"EOF") { // 处理服务端关闭
			fmt.Println("Remote Server close connection.")
			break
		}
		ErrorHandler(err)
		if str != "" {
			fmt.Println("Client received:\n" + str + "from:",add)
		}
	}
}

/*
客户端运行主函数
*/
func cRun(c *Client) {
	err := c.connect() // 连接服务端端口
	if err != nil {
		os.Exit(1)
	}
	go receive(c) 
}

/*
客户端发送消息函数
*/
func runSend(c *Client) {
	scanner := bufio.NewScanner(os.Stdin)
	str := ""
	for scanner.Scan() { // 输入
		input := scanner.Text()
		if input == "exit" { // 输入是exit时，退出程序
			os.Exit(0)
		} else if input == "send"{
			/*
			当输入是send时，将此前在缓冲区的信息发送至server
			eg. 
			Hello World!
			How are you?
			send
			将会发送:
			Hello World!
			How are you?
			*/
			clsStr(&str) //清除开头多余回车与空格
			fmt.Println(str)
			err := c.sendInfo(str)
			if err != nil {
				fmt.Println("Send Error!")
			} else {
				fmt.Println("Client sended:\n" + str)
			}
			str = ""
		} else {
			str += (input + "\n")
		}
	}
}