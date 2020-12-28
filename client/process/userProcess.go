package process

import (
	"fmt"
	"encoding/json"
	"encoding/binary"
	"net"
	"app/chatroom/common/message"
	"app/chatroom/client/utils"
	"os"
)

type UserProcess struct {
	// ...
}

// 注册
func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close() // 延时关闭连接

	// 2. 准备发送消息给服务器
	var mes message.Message 
	mes.Type = message.RegisterMesType

	// 3. 创建注册消息结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4. 序列化注册消息
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 5. 将data（序列化后）赋给消息
	mes.Data = string(data)

	// 6. 序列化消息（mes：将要发送到服务器的消息）
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 创建transfer
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送消息错误，err=", err)
		return
	}

	mes, err = tf.ReadPkg() // mes就是RegisterResMes
	if err != nil {
		fmt.Println("readPkg() err=", err)
		return
	}

	// 将mes的Data反序列化成 RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功！")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return
}

// 登录
func (this *UserProcess) Login(userId int, userPwd string) (err error) {

	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	// 延时关闭连接
	defer conn.Close()

	// 2. 通过 conn 发送消息给服务器
	var mes message.Message 
	mes.Type = message.LoginMesType
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	// 序列化登录消息
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	// 序列化消息（将要发送到服务器的消息）
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 先发送data长度到服务器
	// 获取data长度，转成一个可以表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buf) fail:", err)
		return
	}

	// fmt.Printf("客户端发送消息的长度=%d, 内容=%s \n", len(data), string(data))

	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail:", err)
		return
	}

	// 休眠20s
	// time.Sleep(20 * time.Second)
	// fmt.Println("休眠20s")

	// 处理服务器返回的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg() err=", err)
		return
	}

	// 将mes的Data反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		// 初始化 curUser
		curUser.Conn = conn
		curUser.UserId = userId
		curUser.UserStatus = message.Online

		// 显示当前在线用户列表
		for _, v := range loginResMes.UserIds {
			// 初始化客户端 onlineUsers
			user := &message.User{
				UserId: v,
				UserStatus: message.Online,
			}
			onlineUsers[v] = user
		}

		fmt.Print("\n\n")

		// 启动一个协程：保持和服务器端通信
		// 如果服务器端有数据推送给客户端，则接收并显示
		go serverProcessMes(conn)

		// 显示登录成功后的菜单
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}

	return
}