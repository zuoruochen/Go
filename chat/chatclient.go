//udp传输客户端，还未做重传处理
//服务器返回消息直接打印，报文无格式，未做处理
//TODO：实现可靠性，添加重传机制,消息限制大小，超过大小分多条处理
//version 1.0
//by zrc

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

//包头映射
var reflectString = map[string]string{
	"连接":   "connect:",
	"上线":   "online:",
	"隐身":   "offline:",
	"聊天":   "chat:",
	"在线用户": "get:",
	"退出":   "exit:",
	"心跳":   "heartbeat:",
}

//服务器监听端口
const SERVERPORT = 10000

//用户监听端口
//const CLIENTPORT = 10002

//数据缓冲区大小
const BUFFSIZE = 1024

//全局缓冲区
var buff = make([]byte, BUFFSIZE)

//用户名
var userName string


//错误处理函数
func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

//发送消息给服务器,聊天时回车代表一条信息结束
func sendMessage(udpConn *net.UDPConn, command string) {
//	fmt.Println(command)
	switch command {
	case "chat:":
		fmt.Println("which one to chat:")
		people := ""
		fmt.Scanf("%s", &people)
//		udpConn.Write([]byte(reflectString["聊天"] + people))
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		handleError(err)
		line = "chat:"+userName+":"+people+":"+line
		udpConn.Write([]byte(line))
	default:
		fmt.Println(command)
		udpConn.Write([]byte(command+userName))
	}
}

//处理服务器消息,简单的打印,不区分command类型
func handleMessage(udpListen *net.UDPConn) {
//	fmt.Println("-----handleMessage Begin!-----")
	for {
		n, _, err := udpListen.ReadFromUDP(buff)
//		fmt.Println("handle message!")
		handleError(err)
		if n > 0 {
			fmt.Printf("%s\n", string(buff[:n]))
		}
	}
}

//每秒发送一次心跳包

func heartBeat(udpConn *net.UDPConn) {
	for {
		udpConn.Write([]byte("heartbeat:"+userName))
		time.Sleep(time.Second)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s hostip clientPort\n", os.Args[0])
		os.Exit(1)
	}

	hostIP := os.Args[1]
	clientPort := os.Args[2]
	udpServerAddr, err := net.ResolveUDPAddr("udp", hostIP+":"+strconv.Itoa(SERVERPORT))
	handleError(err)
	udpConnect, err := net.DialUDP("udp", nil, udpServerAddr)
	handleError(err)
	defer udpConnect.Close()

	udpLocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+clientPort)
	handleError(err)
	udpListen, err := net.ListenUDP("udp",udpLocalAddr)
	handleError(err)
	defer udpListen.Close()

	fmt.Println("输入用户名:")
	fmt.Scanf("%s", &userName)

	udpConnect.Write([]byte(reflectString["连接"] + userName + ":" + clientPort))

	go heartBeat(udpConnect)
	go handleMessage(udpListen)

	var command string
	help := `命令(command):
			"上线":				"online"
			"隐身":				"offline"
			"聊天":				"chat"
			"在线用户":			"get"
			"退出":				"exit"`
	fmt.Println(help)

	for {
//		fmt.Printf("[%s@chat client]$",userName)
		fmt.Scanf("%s", &command)
		switch command {
		case "online", "offline", "chat", "get":
			sendMessage(udpConnect, command+":")
		case "exit":
			sendMessage(udpConnect, "exit:")
			os.Exit(0)
		default:
			fmt.Println("输入错误！")
		}
	}
}
