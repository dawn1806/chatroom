package main

import (
	"fmt"
	"net"
	"app/chatroom/common/message"
	"app/chatroom/server/utils"
	"app/chatroom/server/process"
	"io"
)

type Precessor struct {
	Conn net.Conn
}

// 编写一个serverProcessMes函数
// 功能：根据客户端发送消息类型不同，决定调用那个函数来处理
func (this *Precessor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
		case message.LoginMesType:
			// 处理登录，由userPrecess来完成
			// 创建实例
			up := &process.UserProcess{
				Conn: this.Conn,
			}
			err = up.ServerProcessLogin(mes)
		case message.RegisterMesType:
			// 处理注册
			up := &process.UserProcess{
				Conn: this.Conn,
			}
			err = up.ServerProcessRegister(mes)
		case message.SmsMesType:
			// 转发群聊消息
			smsProcess := &process.SmsProcess{}
			smsProcess.SendGroupMes(mes)
		default:
			fmt.Println("消息类型不存在，无法处理...")
	}

	return
}

func (this *Precessor) process2() (err error) {
	// 读客户端发送的消息
	for {
		// 读取数据包，通过工具结构体方法完成
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务也退出...")
			} else {
				fmt.Println("readPkg err=", err)
			}
			return err
		}

		fmt.Println("mes=", mes)

		// 处理客户端发来的消息
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}