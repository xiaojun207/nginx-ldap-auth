package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xiaojun207/nginx-ldap-auth/g"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	uname := this.GetSession("uname")
	if uname == nil {
		this.Ctx.Redirect(302, "/auth/login")
		return
	}
	this.Ctx.Output.Body([]byte("nginx-ldap-auth, version " + g.VERSION + ", this user:" + uname.(string)))
}
