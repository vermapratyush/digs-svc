// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"digs/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//REST
	beego.Router("/v1/login", &controllers.LoginController{})
	beego.Router("/v1/logout", &controllers.LogoutController{})
	beego.Router("/v1/debug", &controllers.DebugController{})

	//WS
	beego.Router("/ws/v1/messenger", &controllers.WSMessengerController{})
}
