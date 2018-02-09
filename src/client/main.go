package main

import (
	"encoding/binary"
	"net"
	"fmt"
)

type tcpClient interface{
	send_msg(data []byte)
	recv_msg()
	colse()
}

type client_tcp struct {
	conn net.Conn 
}

func main() {

	// 实例化 client
	addr:= `127.0.0.1:3563`
	tcp_client := New_client_tcp(addr)

	// Hello 消息（JSON 格式）  对应游戏服务器 Hello 消息结构体
	data := []byte(`{
		"Hello": {
			"Name": "leaf"
		}
	}`)

	// 发送消息
	tcp_client.send_msg(data)

	// 接收消息
	tcp_client.recv_msg()

	// 释放连接
	tcp_client.colse()
}

// 实例化对象
func New_client_tcp(addr string) tcpClient {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	client_tcp := &client_tcp{
		conn:conn,
	}

	return client_tcp
}

// 发送消息
func (this *client_tcp)send_msg(data []byte){
	// 封包
	m:= this.pack_msg(data)

	// 发送消息
	this.conn.Write(m)
}

// 接收消息
func (this *client_tcp)recv_msg(){

	// 解包
	msg_data:= this.unpack_msg()

	fmt.Println(string(msg_data))
}

// 关闭连接
func (this *client_tcp) colse(){
	 this.conn.Close()
}


// 封包
func (this *client_tcp) pack_msg(data []byte) []byte {
	// len + data
	m := make([]byte, 2+len(data))

	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(len(data)))

	copy(m[2:], data)

	return m
}

// 解包
func(this *client_tcp) unpack_msg() (msg_data []byte) {

	// 接收消息 获取消息长度
	msg_head := make([]byte, 2)
	n, err := this.conn.Read(msg_head)
	if err != nil  {
		fmt.Println("读消息头部异常:", err)
	}
	// 默认使用大端序
	msg_len:= binary.BigEndian.Uint16(msg_head[0:n])

	// 读取消息内容
	msg_data = make([]byte, msg_len)
	n, err = this.conn.Read(msg_data)
	if err != nil {
		fmt.Println("读取服务器数据异常:", err)
	}
	return
}