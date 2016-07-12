package socket

import (
	"runtime/debug"
	"digs/logger"
)

func DeadSocketWrite(uid string) {
	if r := recover(); r != nil {
		LeaveNode(uid)
		logger.Critical("PossiblyDeadSocketWrite|FaultyUID=", uid, "|Recovering from panic in MulticastMessage, %v", r, "|Stacktrace=", string(debug.Stack()))
	}
}