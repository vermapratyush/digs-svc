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
	"digs/bots"
)

func init() {
	//REST
	beego.Router("/v1/login", &controllers.LoginController{})
	beego.Router("/v1/logout", &controllers.LogoutController{})
	beego.Router("/v1/notification", &controllers.NotificationController{})
	beego.Router("/v1/unread", &controllers.UnreadController{})
	beego.Router("/v1/people", &controllers.PeopleController{})
	beego.Router("/v1/group", &controllers.GroupController{})
	beego.Router("/v1/feed", &controllers.FeedController{})
	beego.Router("/v1/setting", &controllers.SettingController{})
	beego.Router("/v1/abuse", &controllers.AbuseController{})
	beego.Router("/v1/debug", &controllers.DebugController{})

	//WS
	beego.Router("/ws/v1/messenger", &controllers.WSMessengerController{})

	//BOTS
	beego.Router("/bots/meetup", &bots.MeetupBotController{})
}
