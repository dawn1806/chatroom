package process

import (
	"fmt"
	"app/chatroom/common/message"
	"encoding/json"
)

// 输出群发消息
func outputGroupMes(mes *message.Message) {

	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err.Error())
		return
	}

	// 显示消息
	info := fmt.Sprintf("用户Id: %d 对大家说：%s", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}