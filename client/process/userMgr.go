package process

import (
	"fmt"
	"app/chatroom/common/message"
	"app/chatroom/client/model"
)

// 客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
// 用户登录成功，再对curUser进行初始化
var curUser model.CurUser

// 在客户端显示在线用户列表
func showOnlineUser() {
	fmt.Println("当前在线用户列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户Id：", id)
	}
}

// 更新在线用户状态
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	// 先从在线用户列表查询，看有没有该用户，没有则创建，有则更新
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}

	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	showOnlineUser()
}