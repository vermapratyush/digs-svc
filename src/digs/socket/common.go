package socket

import (
	"github.com/astaxie/beego"
	"runtime/debug"
)

func DeadSocketWrite(peer Peer) {
	if r := recover(); r != nil {
		LeaveNode(peer.UID)
		beego.Critical("PossiblyDeadSocketWrite|FaultyUID=", peer.UID, "|Recovering from panic in MulticastMessage", r, "|Stacktrace=", string(debug.Stack()))
	}
}