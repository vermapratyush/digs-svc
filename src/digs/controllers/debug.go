package controllers

import "digs/socket"

type DebugController struct {
	HttpBaseController
}

func (this *DebugController) Get() {
	debug := make(map[string]interface{})
	lookUp := socket.GetCopy()
	debug["totalActiveClient"] = len(lookUp)
	for k, _ := range(lookUp) {
		debug[k] = 1
	}
	this.Serve200(debug)
}
