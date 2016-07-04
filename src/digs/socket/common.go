package socket

import (
	"github.com/astaxie/beego"
	"runtime/debug"
)

func DeadSocketWrite(uid string) {
	if r := recover(); r != nil {
		LeaveNode(uid)
		beego.Critical("PossiblyDeadSocketWrite|FaultyUID=", uid, "|Recovering from panic in MulticastMessage", r, "|Stacktrace=", string(debug.Stack()))
	}
}