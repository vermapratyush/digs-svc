package socket

import (
	"runtime/debug"
	"digs/logger"
	"digs/domain"
	"digs/models"
)

func DeadSocketWrite(uid string) {
	if r := recover(); r != nil {
		LeaveNode(uid)
		logger.Critical("PossiblyDeadSocketWrite|FaultyUID=", uid, "|Recovering from panic in MulticastMessage, %v", r, "|Stacktrace=", string(debug.Stack()))
	}
}

func AddToFeed(toUID, toGID string, msg *domain.MessageSendRequest) {
	if toGID != "" {
		models.AddToGroupFeed(toGID, msg.MID)
	} else {
		models.AddToUserFeed(toUID, msg.MID)
	}
}