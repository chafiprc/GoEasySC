package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	port int // 端口
	conn net.Conn // 连接
	isConnected bool // 是否连接
	remoteAddress net.Addr // 服务端地址
}

/*
发送消息
*/
func (c *Client) sendInfo(str string) error {
	if !c.isConnected {
		return fmt.Errorf("not connected")
	}

	_, err := c.conn.Write([]byte(str))

	if err != nil {
		return err
	}
	return nil
}

/*
接受信息
*/
func (c *Client) receiveInfo() (string,net.Addr,error) {
	if !c.isConnected {
		return "", c.remoteAddress, fmt.Errorf("not connected")
	}
	
	buffer := make([]byte,1024)

	num, err:= c.conn.Read(buffer)

	if err != nil {
		if strings.Contains(err.Error(),"An existing connection was forcibly closed by the remote host.") {
			c.isConnected = false
			fmt.Println("Remote Server closed.")
		}
		return "",c.remoteAddress,err
	}
	data := string(buffer[:num])

	return data, c.remoteAddress, nil
}

/*
初始化
*/
func (c *Client) init() error {
	c.isConnected = false
	fmt.Println("Input port:")
	_,err := fmt.Scan(&c.port)
	if err != nil {
		return err
	}
	return nil
}

/*
连接服务端
*/
func (c *Client) connect() error {
	conn, err := net.Dial("tcp","127.0.0.1:"+strconv.Itoa(c.port))
	if err != nil {
		fmt.Println("Error Connection!",err)
		c.isConnected = false
		return err
	}
	c.remoteAddress = conn.RemoteAddr()
	c.conn = conn
	c.isConnected = true

	fmt.Println("Succeeded in connecting to server, address:",c.conn.RemoteAddr())
	return nil
}

/*
关闭客户端
*/
func (c *Client) close() error {
	if !c.isConnected {
		return nil
	}
	err := c.conn.Close()
	if err != nil {
		return err
	}
	c.isConnected = false
	return nil
}