package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// 注释同CLient
	s := new(Server)
	err := s.init()
	ErrorHandler(err)
	defer s.close()
	go sRun(s)
	runSend(s)
}

func receive(s *Server) {
	for {
		str, add, err := s.receiveInfo()
		if err != nil {
			if strings.Contains(err.Error(),"not connected") {
				fmt.Println("Remote connected was closed. Or not connected to remote host.")
				break
			}
		}
		ErrorHandler(err)
		if str != "" {
			fmt.Println("Server received:\n" + str + "from:",add)
			s.messageHandler(str)
		}
	}
}

func sRun(s *Server){
	for { // 用于远程客户端关闭后还能继续连接
		if !s.isConnected {
			s.connect()
			go receive(s)
		}
		time.Sleep(1*time.Second) // 1秒刷新
	}
}

func runSend(s *Server) {
	scanner := bufio.NewScanner(os.Stdin)
	str := ""
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			os.Exit(0)
		} else if input == "send"{
			for str[0] == '\n' || str[0] == ' '{
				str = str[1:]
			}
			err := s.sendInfo(str)
			if err != nil {
				fmt.Println("Send Error!")
			} else {
				fmt.Println("Server sended:\n" + str)
			}
			str = ""
		} else {
			str += (input + "\n")
		}
	}
}