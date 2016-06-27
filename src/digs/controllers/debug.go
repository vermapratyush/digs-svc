package controllers

import "digs/socket"

type DebugController struct {
	HttpBaseController
}

func (this *DebugController) Get() {
	debug := make(map[string]interface{})
	debug["totalActiveClient"] = len(socket.LookUp)
	for k, _ := range(socket.LookUp) {
		debug[k] = 1
		debug["push-" + k] = socket.LookUp[k].PushNotificationEnabled
	}
	this.Serve200(debug)
}
