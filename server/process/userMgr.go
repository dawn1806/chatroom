package process

import (
	"fmt"
)

// userMgr有且只有一个，在很多地方都用到，所以声明为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 初始化
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 增加onlineUser
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

// 删除onlineUser
func (this *UserMgr) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

// 查询onlineUser（全部）
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

// 根据Id返回onlineUser
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	// 利用带检测的方式，从map中取出一个值
	up, ok := this.onlineUsers[userId]
	if !ok {
		// 这里说明要查找的这个用户 当前不在线
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}