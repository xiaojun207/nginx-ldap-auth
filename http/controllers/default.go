package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	uname := this.GetSession("uname")

	data := map[string]interface{}{}

	if uname == nil {
		data = map[string]interface{}{
			"code": "100100",
			"msg":  "未登录",
		}
	} else {
		data = map[string]interface{}{
			"code": "100200",
			"msg":  "成功",
			"data": uname,
		}
	}

	this.Data["json"] = data
	this.ServeJSON()
	//this.Ctx.Output.Body([]byte("nginx-ldap-auth, version " + g.VERSION + ", this user:" + uname.(string)))
}
