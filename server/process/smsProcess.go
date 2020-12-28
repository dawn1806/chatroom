package process

import (
	"net"
	"fmt"
	"app/chatroom/client/utils"
	"app/chatroom/common/message"
	"encoding/json"
)

type SmsProcess struct {
	// 暂时不需要字段
}

func (this *SmsProcess) SendGroupMes(mes *message.Message) {

	// 取出消息内容
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	// 这里只是转发消息，消息内容在这里不需要知道
	// 所以将 mes 序列化后发送即可
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 遍历在线用户列表，分别发送消息
	for id, up := range userMgr.onlineUsers {
		// 这里过滤掉自己，不要把消息发送给自己
		if id == smsMes.UserId {
			continue
		}
		this.sendMes(data, up.Conn)
	}
}

func (this *SmsProcess) sendMes(data []byte, conn net.Conn) {

	// 创建transfer进行发送
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}