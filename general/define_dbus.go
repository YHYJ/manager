/*
File: define_dbus.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-13 13:10:41

Description: D-Bus总线通信
*/

package general

import (
	"os"

	"github.com/godbus/dbus/v5"
)

// PolicyKitAuthentication 使用D-Bus进行身份认证
func PolicyKitAuthentication() {
	// PolicyKit 认证，通过 D-Bus 总线通信
	// 1. 连接到系统总线 SystemBus
	systemBus, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}

	// 2. 获取 PolicyKit 授权代理对象
	// Object 方法的第一个参数是 Bus name，第二个参数是 Object path
	policykitObject := systemBus.Object("org.freedesktop.PolicyKit1", "/org/freedesktop/PolicyKit1/Authority")
	if policykitObject == nil {
		panic("Unable to perform Polkit authentication")
	}

	// 3. 构建认证方法的参数('subject', 'action_id', 'details', 'flags', 'cancellation_id')
	// 'subject' 需要新建一个如下结构体，内容是本进程的信息，包含进程ID和启动时间
	type Subject struct {
		Type       string
		Identifier map[string]dbus.Variant
	}
	subject := Subject{
		Type: "unix-process",
		Identifier: map[string]dbus.Variant{
			"pid":        dbus.MakeVariant(uint32(os.Getpid())),
			"start-time": dbus.MakeVariant(uint64(0)),
		},
	}
	// 'action_id' 是一个字符串，用于标识本次认证的目的
	actionID := "org.freedesktop.policykit.exec"
	// 'details' 是一个map，用于描述本次认证的详细信息
	// details := map[string]dbus.Variant{
	details := map[string]string{}
	// 'flags'是一个uint32，用于标识本次认证的类型
	flags := uint32(1)
	// 'cancellation_id' 是一个字符串，作为本次认证取消的标识
	cancellationID := ""

	// 4. 创建认证结果的接收对象
	var result dbus.Variant

	// 5. 调用 CheckAuthorization 方法进行认证
	err = policykitObject.Call(
		"org.freedesktop.PolicyKit1.Authority.CheckAuthorization", 0,
		subject, actionID, details, flags, cancellationID).Store(&result)
	if err != nil {
		panic(err)
	}

	// 6. 解析认证结果
	// 认证成功时，result 的值是一个结构体，包含认证成功的信息
	value := result.Value().([]interface{})[0]
	if value.(bool) != true {
		panic("Error executing command as another user: Not authorized")
	}

	// 7. 认证成功，执行需要认证的操作
}
