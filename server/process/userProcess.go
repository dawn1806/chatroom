package process

import (
	"fmt"
	"net"
	"app/chatroom/common/message"
	"app/chatroom/server/utils"
	"app/chatroom/server/model"
	"encoding/json"
)

type UserProcess struct {
	Conn net.Conn
	UserId int // 为了表示Conn是哪个用户的
}

// userId代表的用户通知其他在线用户：我上线了
func (this *UserProcess) NotifyOtherOnlineUser(userId int) {
	// 遍历 onlineUsers 进行通知，消息结构是 NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {

		// 过滤掉自己
		if id == userId {
			continue
		}

		// 开始一对一通知: up代表当前要通知的用户
		up.NotifyOneByOneOnlineUser(userId)
	}
}

func (this *UserProcess) NotifyOneByOneOnlineUser(userId int) {
	// 组装通知消息
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.Online

	// 序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyOneByOneOnlineUser err=", err)
		return
	}
}

// 此函数：处理注册用户
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 从mes中取出mes.Data， 反序列化成 RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	// 声明resMes，返回到客户端
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	// 声明 registerResMes ，存储返回的消息
	var registerResMes message.RegisterResMes

	// 在redis中完成注册
	err = model.GlobalUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 400
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 509
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	// 进行赋值
	resMes.Data = string(data)

	// 对resMes进行序列化，准备发送回客户端
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	// 发送（使用自定义工具函数 writePkg）
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

// 此函数：处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 核心代码：
	// 1. 先从mes中取出mes.Data， 反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	// 声明resMes，返回到客户端
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	// 声明loginResMes，存储返回的消息
	var loginResMes message.LoginResMes

	// 从redis中验证用户
	user, err := model.GlobalUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
		
	} else {
		loginResMes.Code = 200
		fmt.Printf("用户【%v】登录成功", user)

		// 用户登录成功后，将用户放入 UserMgr(维护了在线用户列表)
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)

		// 通知其他在线用户：我上线了
		this.NotifyOtherOnlineUser(loginMes.UserId)

		// 将 在线用户Id 放入 loginResMes 返回客户端
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UserIds = append(loginResMes.UserIds, id)
		}
	}

	// 如果用户Id=100，密码=123456，认为合法，否则不合法
	// if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	// 	loginResMes.Code = 200
	// } else {
	// 	loginResMes.Code = 500 // 500状态码，表示该用户不存在
	// 	loginResMes.Error = "该用户不存在，请先注册"
	// }

	// 将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	// 进行赋值
	resMes.Data = string(data)

	// 对resMes进行序列化，准备发送回客户端
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	// 发送（使用自定义工具函数 writePkg）
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}