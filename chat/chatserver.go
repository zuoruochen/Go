//单gorutine处理信息模型，缺乏锁机制
//TODO：改为多goruntine，添加锁机制
//version 1.0
//by zrc
package main

import (
//	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

//服务器端保留的用户信息
type user struct {
	userName       string
	userAddr       *net.UDPAddr
	userListenConn *net.UDPConn //用户接受消息套接字
	active			bool		//用户是否登入服务器，用户掉线将改为false
}

//服务器监听端口
const LISTENPORT = 10000

//数据缓冲区大小
const BUFFSIZE = 1024

//全局缓冲区
var buff = make([]byte, BUFFSIZE)

//活跃用户列表,登入了服务器未主动退出的用户列表,key为username
var activeUser map[string]user
//在线用户列表,key为username
var onlineUser map[string]bool

var heartbeatCheck map[string]bool

//错误处理函数
func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}


func onlineCheck() {
	for {
		time.Sleep(5*time.Second)
		for key,value := range heartbeatCheck {
			if value {
				heartbeatCheck[key] = false
			}else {
				onlineUser[key] = false
				temp := activeUser[key]
				temp.active = false
				activeUser[key] = temp
				fmt.Println(key+" 掉线了！")
			}
		}
	}
}
/*
 *connect:username:port
 *online:username
 *offline:username
 *get:
 *chat:username:usertochat:message
 *exit:username
 *heartbeat:username
*/

func parseMessage(buff []byte,len int) []string {
	analMsg := make([]string, 0)
	strNow := ""
	for i := 0; i < len; i++ {
		if string(buff[i:i + 1]) == ":" {
			analMsg = append(analMsg, strNow)
			strNow = ""
		} else {
			strNow += string(buff[i:i + 1])
		}
	}
	analMsg = append(analMsg, strNow)
	return analMsg
}

func handleMessage(udpListen *net.UDPConn) {
//	fmt.Println("----handleMessage begin----")
//	n,_,err := udpListen.ReadFromUDP(buff)
	n,addr,err := udpListen.ReadFromUDP(buff)
	handleError(err)
//	fmt.Println("start handling message!")
	if n > 0 {
		msgList := parseMessage(buff,n)
//		fmt.Print("[服务器信息]：")
		switch msgList[0] {
			case "connect":
				userName := msgList[1]
				userListenPort := msgList[2]
//				fmt.Printf("%v\n",addr.IP)
				userIP := parseMessage([]byte(addr.String()), len(addr.String()))
				fmt.Printf("%s 登入服务器, IP地址:%s,端口: %s\n", userName,userIP[0],userListenPort)
				userAddr, err := net.ResolveUDPAddr("udp", userIP[0]+":"+userListenPort)
				handleError(err)
				userListenConn, err := net.DialUDP("udp", nil, userAddr)
				handleError(err)
//				activeUserList = append(activeList,user{userName,userAddr,userListenConn,nil})
				activeUser[userName] = user{userName,userAddr,userListenConn,true}
				userListenConn.Write([]byte("connect done!"))
//				onlineUser[userName] = false

			case "online":
				fmt.Printf("%s 上线了!\n",msgList[1])
				userName := msgList[1]
				onlineUser[userName] = true
//				for key,value := range onlineUser {
//					fmt.Println(key,value)
//				}

			case "offline":
				fmt.Printf("%s 下线了！\n",msgList[1])
				//修改onlineList
				userName := msgList[1]
				delete(onlineUser,userName)
			case "chat":
				fmt.Println("chat")
				//TODO：msgList[1]发送msgList[3]给msgList[2]
				userName := msgList[1]
				userToChatName := msgList[2]
				message := msgList[3]
				fmt.Println(message)
				userListenConn := activeUser[userToChatName].userListenConn
				userListenConn.Write([]byte(userName+" said : " + message))
			case "get":
				fmt.Println("get")
				//发送onlineUser
				userName := msgList[1]
				userListenConn := activeUser[userName].userListenConn
				activeList := ""
				for key,value:= range onlineUser {
//					fmt.Println(key)
					if value {
						activeList = activeList + key + "\n"
					}
				}
				fmt.Println(activeList)
				userListenConn.Write([]byte(activeList))
			case "exit":
				fmt.Printf("%s 退出了服务器\n",msgList[1])
				//修改activeUser和onlineUser，关闭套接字
				userName := msgList[1]
				delete(onlineUser,userName)
				activeUser[userName].userListenConn.Close()
				delete(activeUser,userName)
				delete(heartbeatCheck,userName)

			case "heartbeat":
				//fmt.Println("heartbeat")
				//查看msgList[1]是否在activeUser中,如果不在则将其加入
				userName := msgList[1]
				heartbeatCheck[userName] = true
				temp := activeUser[userName]
				temp.active= true
				activeUser[userName] = temp
				v,find := onlineUser[userName]
				if find {
					if !v {
						onlineUser[userName] = true
					}
				}
			default:
				fmt.Println("未知信息！")
		}
	}
}

func main() {
	udpListenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(LISTENPORT))
	handleError(err)
	udpListen, err := net.ListenUDP("udp",udpListenAddr)
	handleError(err)
	defer udpListen.Close()

	fmt.Println("开始监听：")
	activeUser = make(map[string]user)
	onlineUser = make(map[string]bool)
	heartbeatCheck = make(map[string]bool)
	go onlineCheck()

	for {
		handleMessage(udpListen)
	}
}
