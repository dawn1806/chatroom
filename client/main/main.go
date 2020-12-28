package main

import (
	"fmt"
	"os"
	"app/chatroom/client/process"
)

// 全局变量：用户信息
var userId int
var userPwd string 
var userName string

func main () {

	// 定义key 接收用户的选择
	var key int
	// 判断是否继续显示菜单
	// var loop = true

	for true {
		fmt.Println("---------欢迎登陆多人聊天室--------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("请选择（1-3）：")

		fmt.Scanf("%d\n", &key)
		switch key {
			case 1:
				fmt.Println("登陆聊天室")
				fmt.Println("输入用户Id：")
				fmt.Scanf("%d\n", &userId)
				fmt.Println("输入密码：")
				fmt.Scanf("%s\n", &userPwd)
				// 登录
				up := &process.UserProcess{}
				up.Login(userId, userPwd)

			case 2:
				fmt.Println("注册用户")
				fmt.Println("输入用户Id：")
				fmt.Scanf("%d\n", &userId)
				fmt.Println("输入密码：")
				fmt.Scanf("%s\n", &userPwd)
				fmt.Println("输入用户名（昵称）：")
				fmt.Scanf("%s\n", &userName)
				// 注册
				up := &process.UserProcess{}
				up.Register(userId, userPwd, userName)

			case 3:
				fmt.Println("退出系统")
				os.Exit(0)
			default:
				fmt.Println("你的输入有误，请重新输入")
		}
	}

	// 根据用户的输入，显示新的信息
	// if key == 1 {
	// 	// 用户要登陆了
	// 	fmt.Println("输入用户Id：")
	// 	fmt.Scanf("%d\n", &userId)
	// 	fmt.Println("输入密码：")
	// 	fmt.Scanf("%s\n", &userPwd)
	// 	// 将登陆的业务逻辑写入另一个文件 login.go
	// 	login(userId, userPwd)

	// } else if key == 2 {
	// 	fmt.Println("这里是注册用户")
	// }	
}