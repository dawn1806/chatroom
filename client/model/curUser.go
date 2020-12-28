package model

import (
	"net"
	"app/chatroom/common/message"
)

// 在userMgr中将其定义为全局变量
type CurUser struct {
	Conn net.Conn
	message.User
}