package main

import (
	"fmt"
	"encoding/json"
	"encoding/binary"
	"net"
	"app/chatroom/common/message"
	_ "time"
)

func login (userId int, userPwd string) (err error) {

	// fmt.Println("用户Id=%d\n", userId)
	// fmt.Println("用户密码=%s\n", userPwd)

	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	// 延时关闭连接
	defer conn.Close()

	// 2. 通过 conn 发送消息给服务器
	var mes message.Message 
	mes.Type = message.LoginMesType
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	// 序列化登录消息
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	// 序列化消息（将要发送到服务器的消息）
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 先发送data长度到服务器
	// 获取data长度，转成一个可以表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buf) fail:", err)
		return
	}

	// fmt.Printf("客户端发送消息的长度=%d, 内容=%s \n", len(data), string(data))

	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail:", err)
		return
	}

	// 休眠20s
	// time.Sleep(20 * time.Second)
	// fmt.Println("休眠20s")

	// 处理服务器返回的消息
	mes, err = readPkg(conn)
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	// 将mes的Data反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登录成功")
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}

	return
}