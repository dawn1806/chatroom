package message

const (
	LoginMesType			= "LoginMes"
	LoginResMesType			= "LoginResMes"
	RegisterMesType 		= "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType				= "SmsMes"
)

// 用户状态常量
const (
	Online = iota
	Offline
	Busy
)

type Message struct {
	Type string `json:"type"`// 消息类型
	Data string `json:"data"`// 消息数据
}

// 定义具体的消息
type LoginMes struct {
	UserId int `json:"userId"`// 用户Id
	UserPwd string `json:"userPwd"`// 用户密码
	UserName string `json:"userName"`// 用户名
}

type LoginResMes struct {
	Code int `json:"code"`// 状态码 500表示未注册  200表示登录成功
	UserIds []int `json:"userIds"` // 在线用户Id集合，返回客户端，让客户端知道当前在线用户有哪些
	Error string `json:"error"`// 错误消息
}

type RegisterMes struct {
	User User `json:"user"`
}

type RegisterResMes struct {
	Code int `json:"code"`// 状态码 400表示用户已存在  200表示注册成功
	Error string `json:"error"`// 错误消息
}

// 用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

// 发送的消息
type SmsMes struct {
	Content string `json:"content"` // 消息内容
	User // 匿名结构体 继承
}