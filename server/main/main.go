package main

import (
	"fmt"
	"net"
	"time"
	"app/chatroom/server/model"
)

// 处理和客户端的通讯
func process1 (conn net.Conn) {
	// 延时关闭连接
	defer conn.Close()
	
	// ...
	processor := &Precessor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务端通讯的协程错误，err=", err)
		return
	}
}

// 初始化UserDao
func initUserDao() {
	// pool：在redis.go定义的全局变量，在initPool中初始化
	// 此函数需要在initPool后面执行
	model.GlobalUserDao = model.NewUserDao(pool)
}

func main () {

	// 当服务器启动时，就初始化redis连接池
	initPool("127.0.0.1:6379", 16, 0, 300 * time.Second)
	initUserDao()

	// 提示信息
	fmt.Println("服务器[新的结构]在 8889 端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}

	// 监听成功后，延时关闭listen
	defer listen.Close()

	// 监听成功后，就等待客户端来连接服务器
	for {
		fmt.Println("等待客户端来连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
			// 这里不用return，某一个客户端连接失败，可以继续等待其他客户端连接
		}

		// 连接成功，启动一个协程和客户端保持通讯
		go process1(conn)
	}
}