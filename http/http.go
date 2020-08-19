package http

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/xiaojun207/nginx-ldap-auth/g"
	"github.com/xiaojun207/nginx-ldap-auth/http/controllers"
)

func Start() {
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "sessionID"
	beego.BConfig.WebConfig.EnableXSRF = true
	beego.BConfig.WebConfig.XSRFKey = "ed769515656b704dee92d77e28663147"
	beego.BConfig.WebConfig.XSRFExpire = 3600

	beego.SetStaticPath("/auth/static/", "views/static")

	if !g.Config().Http.Debug {
		logs.SetLevel(logs.LevelInformational)
	}
	ConfigRouters()
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run(g.Config().Http.Listen)
}
