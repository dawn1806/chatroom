package process

import (
	"fmt"
	"os"
	"net"
	"app/chatroom/client/utils"
	"app/chatroom/common/message"
	"encoding/json"
)

// 显示登录成功后的界面
func ShowMenu() {

	fmt.Println("---------恭喜xxx登录成功--------")
	fmt.Println("\t\t\t 1 显示在线用户列表")
	fmt.Println("\t\t\t 2 发送消息")
	fmt.Println("\t\t\t 3 信息列表")
	fmt.Println("\t\t\t 4 退出系统")
	fmt.Println("请选择（1-4）：")

	var key int
	var content string

	// 因为发送消息比较频繁，所以 smsProcess 实例创建在switch外面
	smsProcess := &SmsProcess{}

	fmt.Scanf("%d\n", &key)
	switch key {
		case 1:
			// 显示在线用户列表
			showOnlineUser()
		case 2:
			fmt.Println("你想对大家说什么：")
			fmt.Scanf("%s\n", &content)
			smsProcess.SendGroupMes(content)
		case 3:
			fmt.Println("信息列表")
		case 4:
			fmt.Println("你选择退出了系统")
			os.Exit(0)
		default:
			fmt.Println("输入有误...")
	}
}

// 和服务器端保持通信
func serverProcessMes(conn net.Conn) {
	// 创建一个transfer实例，不停的读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务端发来的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() err=", err)
			return
		}

		// 读到消息，进行下一步处理
		switch mes.Type {
			case message.NotifyUserStatusMesType:
				// 有人上线了
				// 1. 取出消息 NotifyUserStatusMes
				var notifyUserStatusMes message.NotifyUserStatusMes
				json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
				// 2. 把消息保存到客户端维护的 map[int]User 中
				updateUserStatus(&notifyUserStatusMes)
			case message.SmsMesType:
				// 有人群发消息
				outputGroupMes(&mes)
			default:
				fmt.Println("服务器返回了未知消息类型")
		}
	}
}