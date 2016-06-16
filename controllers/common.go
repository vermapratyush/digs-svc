package controllers

import (
)
import (
	"digs/domain"
)

type BaseController struct {
	Data *map[interface{}]interface{}
	SetStatus func (status int)
	Serve func(encoding ...bool)
}

func (this *BaseController) Serve500(err error) {
	(*this.Data)["json"] = domain.ErrorResponse{
		StatusCode:500,
		Message:err.Error(),
	}
	this.SetStatus(500)
	this.Serve()
	return

}

func (this *BaseController) Serve200(obj interface{}) {
	(*this.Data)["json"] = obj
	this.SetStatus(200)
	this.Serve()
}

func (this *BaseController) Serve304() {
	this.SetStatus(304)
	this.Serve()
}
