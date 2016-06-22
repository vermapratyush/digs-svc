package socket

import "github.com/astaxie/beego"

func DeadSocketWrite() {
	if r := recover(); r != nil {
		beego.Critical("PossiblyDeadSocketWrite| Recovering from panic in MulticastMessage", r)
	}
}