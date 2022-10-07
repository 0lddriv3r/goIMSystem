package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前客户端的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
	}

	//连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error: ", err)
	}

	client.conn = conn

	//返回对象
	return client
}

// 处理server回应的消息，直接显示到客户端标准输出
func (client *Client) DealResponse() {
	/*for {
		buf := make([]byte, 4096)
		client.conn.Read(buf)
		fmt.Println(buf)
	}*/

	//一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1、公聊模式")
	fmt.Println("2、私聊模式")
	fmt.Println("3、修改用户名")
	fmt.Println("0、退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>> 请输入合法范围内的数字 <<<<<")
		return false
	}
}

func (client *Client) PublicChat() {
	//提示用户输入消息
	var chatMsg string
	fmt.Println(">>>>> 请输入聊天内容，exit退出 <<<<<")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		//发送非空消息给服务器
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err: ", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>> 请输入聊天内容，exit退出 <<<<<")
		fmt.Scanln(&chatMsg)
	}
}

// 查询在线用户
func (client *Client) QueryOnlineUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err: ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string

	client.QueryOnlineUsers()

	//提示用户输入消息
	var chatMsg string
	fmt.Println(">>>>> 请输入聊天对象[用户名]，exit退出 <<<<<")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>> 请输入聊天内容，exit退出 <<<<<")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//发送非空消息给服务器
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err: ", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>> 请输入聊天内容，exit退出 <<<<<")
			fmt.Scanln(&chatMsg)
		}

		client.QueryOnlineUsers()
		fmt.Println(">>>>> 请输入聊天对象[用户名]，exit退出 <<<<<")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>> 请输入用户名: <<<<<")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		//根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认: 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认: 8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败 <<<<<")
	}

	go client.DealResponse()

	fmt.Println(">>>>> 连接服务器成功 <<<<<")

	//启动客户端业务
	client.Run()
}
