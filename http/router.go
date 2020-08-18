package http

import (
	"github.com/astaxie/beego"
	"github.com/shanghai-edu/nginx-ldap-auth/http/controllers"
)

func ConfigRouters() {
	beego.Router("/auth", &controllers.MainController{})
	beego.Router("/auth/login", &controllers.LoginController{})
	beego.Router("/auth/logout", &controllers.LogoutController{})
	beego.Router("/auth/auth-proxy", &controllers.AuthProxyController{})
	beego.Router("/auth/api/v1/:control", &controllers.ControlController{})
}
