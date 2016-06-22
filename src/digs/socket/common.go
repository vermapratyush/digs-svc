package socket

import (
	"github.com/astaxie/beego"
	"digs/models"
)

func DeadSocketWrite(userAccount *models.UserAccount) {
	if r := recover(); r != nil {
		LeaveNode(userAccount.UID)
		beego.Critical("PossiblyDeadSocketWrite|FaultyUID=", userAccount.UID, "|Recovering from panic in MulticastMessage", r)
	}
}