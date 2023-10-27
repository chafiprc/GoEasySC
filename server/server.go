package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	port int // 服务端端口
	listener net.Listener // 监听器
	conn net.Conn // 连接
	isConnected bool // 是否连接
	remoteAddress net.Addr // 远程客户端地址
}

/*
接受远程客户端信息，输出接收到的信息以及地址与err
*/
func (s *Server) receiveInfo() (string,net.Addr,error) {
	if !s.isConnected { // 未连接时不可以接受信息
		return "",s.remoteAddress,fmt.Errorf("not connected")
	}
	buffer := make([]byte,1024) // 缓冲区

	num, err:= s.conn.Read(buffer) 

	if err != nil {
		if strings.Contains(err.Error(),"An existing connection was forcibly closed by the remote host.") {
			s.isConnected = false
			fmt.Println("Remote Client closed.")
		}
		return "",s.remoteAddress,err
	}

	data := string(buffer[:num])

	return data, s.remoteAddress, nil
	
}

/*
服务端发送信息函数，传入参数：需要发送的信息
*/
func (s *Server) sendInfo(str string) error {
	if !s.isConnected { // 未连接时不可发送
		return fmt.Errorf("not connected")
	}
	
	_, err := s.conn.Write([]byte(str))
	if err!= nil {
		return err
	}
	fmt.Println("Sended Info:",str)
	return nil
}

/*
服务端初始化函数
*/
func (s *Server) init() error {
	s.isConnected = false // 默认未连接
	fmt.Println("Input port:")
	_, err := fmt.Scan(&s.port) // 输入服务端地址端口
	if err != nil {
		return err
	}
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp","127.0.0.1:"+strconv.Itoa(s.port)) //默认地址是127.0.0.1
	if err!= nil {
		return err
	}

	s.listener = listener
	fmt.Println("Server started successfully, waiting for connection...")
	return nil
}

/*
关闭服务器
*/
func (s *Server) close() error{
	err := s.listener.Close()
	fmt.Println("Close server.")
	return err
}

/*
连接远程客户端
*/
func (s *Server) connect() error{ // 单次连接
	conn, err := s.listener.Accept()
	if err != nil {
		fmt.Println("Connection error.")
		s.isConnected = false
		return err
	}
	s.remoteAddress = conn.RemoteAddr()
	fmt.Println("Succeeded in connection. Remote address:",s.remoteAddress) // 获取远程地址
	s.conn = conn
	s.isConnected = true
	return nil
}

/*
消息处理函数
输出长度为10的字符串数组与err（最多接受10行消息）
*/
func (s *Server) messageHandler(str string) ([10]string,int) {
	strLine := [10]string{}
	lineCnt := 0
	start := 0

	// 进行分行
	for enterIndex:= strings.Index(str[start:],"\n"); enterIndex!=-1&&lineCnt<=9; enterIndex = strings.Index(str[start:],"\n"){		enterIndex += start
		strLine[lineCnt] = str[start:enterIndex]
		start = enterIndex + 1
		lineCnt++ // 统计有效信息有几行
	}

	// 获取method
	methodIndex := strings.Index(strLine[0]," ")
	if methodIndex == -1{
		return strLine, lineCnt-1
	}
	method := strLine[0][:methodIndex]
	if method != "OPTIONS" && method != "SETUP" && method != "PLAY" && method != "TEARDOWN" { // 规避报错
		return strLine, lineCnt-1
	}

	// 获取url
	urlIndex := strings.Index(strLine[0][methodIndex+1:]," ")
	if urlIndex == -1 {
		return strLine, lineCnt-1
	}
	url := strLine[0][methodIndex+1:methodIndex+urlIndex+1]

	// 获取version
	version := strLine[0][methodIndex+urlIndex+2:]

	cseq := "1" //默认序列号
	if lineCnt - 1 >= 2{
		cseq = strLine[1][5:] // 获取序列号
	}
	s.requestHandler(method,url,version,cseq,strLine) // 格式符合，进行request处理，否则不进行处理

	return strLine, lineCnt-1
}


/*
请求处理，输入method，url，version，cseq以及所有接受到的信息

接收到
OPTIONS url version
...
服务端返回
0.5 200 OK
Cseq:x
OPTIONS
SETUP
PLAY
TEARDOWN

接收到
SETUP url version
...
服务端返回
0.5 200 OK
Cseq:x
session_id:

接收到
PLAY url version
...
服务端返回1~99

接收到
TEARDOWN url version
...
服务端关闭客户端连接
*/
func (s *Server) requestHandler(method string,url string, version string, cseq string, info [10]string){
	if method == "OPTIONS" {
		s.sendInfo("0.5 200 OK\nCseq:" + cseq + "\n" + "OPTIONS\nSETUP\nPLAY\nTEARDOWN")
	} else if method == "SETUP" {
		portIndex := strings.Index(s.remoteAddress.String(),":") 
		s.sendInfo("0.5 200 OK\nCseq:" + cseq + "\n" + "session_id:" + s.remoteAddress.String()[portIndex+1:]) //session_id与port一致
	} else if method == "PLAY" {
		for i:=0; i< 100; i++ {
			s.sendInfo(fmt.Sprint(i))
		}
	} else if method == "TEARDOWN" {
		s.conn.Close()
		s.isConnected = false
	}
}