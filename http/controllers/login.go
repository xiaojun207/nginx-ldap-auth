package controllers

import (
	"fmt"
	"time"

	"html/template"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils/captcha"
	"github.com/xiaojun207/nginx-ldap-auth/g"
	"github.com/xiaojun207/nginx-ldap-auth/utils"
)

type LoginController struct {
	beego.Controller
}

func init() {
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/auth/captcha/", store)
	cpt.ChallengeNums = 2
	cpt.StdWidth = 120
	cpt.StdHeight = 40

	createCaptcha := func() template.HTML {
		value, err := cpt.CreateCaptcha()
		if err != nil {
			logs.Error("Create Captcha Error:", err)
			return ""
		}
		// create html
		return template.HTML(fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`+
			`<a class="captcha" href="javascript:">`+
			`<img onclick="this.src=('%s%s?reload='+(new Date()).getTime())" class="captcha-img" src="%s%s">`+
			`</a>`, cpt.FieldIDName, value, cpt.URLPrefix, value, cpt.URLPrefix, value))
	}
	// add to template func map
	beego.AddFuncMap("createcaptcha", createCaptcha)
}

var cpt *captcha.Captcha

func getMsg(loginFailed interface{}) string {
	var msg string
	switch loginFailed {
	case "1":
		msg = "Login Failed: Username Or Password Wrong"
	case "2":
		msg = "Login Failed: User is not Allowed"
	case "3":
		msg = "Login Failed: Captcha Wrong"
	case "4":
		msg = "Login Failed: 尝试次数过多"
	case "5":
		msg = "Login Failed: 尝试过于频繁"
	}
	return msg
}

func (this *LoginController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	logtime := time.Now().Format("02/Jan/2006 03:04:05")
	target := this.Ctx.Input.Header("X-Target")
	getTarget := this.GetString("target")
	if target == "" && getTarget == "" {
		target = "/"
	}
	if getTarget != "" {
		target = getTarget
	}
	this.Data["target"] = target
	loginFailed := this.GetSession("loginFailed")
	if loginFailed != nil {
		this.Data["captcha"] = true
	}
	msg := getMsg(loginFailed)
	this.Data["msg"] = msg
	this.TplName = "template/login.tpl"
	clientIP := this.Ctx.Input.IP()
	DirectIPS := g.Config().Control.IpAcl.Direct
	DenyIPS := g.Config().Control.IpAcl.Deny
	timeDirect := g.Config().Control.TimeAcl.Direct
	timeDeny := g.Config().Control.TimeAcl.Deny
	if utils.IpCheck(clientIP, DenyIPS) {
		logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: IP %s is not allowed", clientIP, logtime, clientIP))
		this.Abort("401")
	}

	if utils.IpCheck(clientIP, DirectIPS) {
		this.SetSession("uname", clientIP)
		logs.Notice(fmt.Sprintf("%s - %s [%s] Login Successed: Direct IP", clientIP, clientIP, logtime))
		this.TplName = "template/direct.tpl"
		return
	}
	if utils.TimeCheck(timeDeny) {
		logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: This Time is not allowed", clientIP, logtime))
		this.Abort("403")
	}
	if utils.TimeCheck(timeDirect) {
		this.SetSession("uname", "timeDirect")
		logs.Notice(fmt.Sprintf("%s - %s [%s] Login Successed: Direct Time", clientIP, "timeDirect", logtime))
		this.TplName = "template/direct.tpl"
		return
	}
}

func (this *LoginController) Post() {
	logtime := time.Now().Format("02/Jan/2006 03:04:05")
	clientIP := this.Ctx.Input.IP()
	this.Ctx.Request.ParseForm()
	username := this.Ctx.Request.Form.Get("username")
	password := this.Ctx.Request.Form.Get("password")
	target := this.Ctx.Request.Form.Get("target")
	loginFailed := this.GetSession("loginFailed")
	if loginFailed != nil {
		if !cpt.VerifyReq(this.Ctx.Request) {
			this.SetSession("loginFailed", "3")
			logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: Captcha Wrong", clientIP, logtime))
			this.Ctx.Redirect(302, fmt.Sprintf("/auth/login?target=%s", target))
			return
		}
	}

	if len(g.Config().Control.AllowUser) > 0 {
		if !utils.In_slice(username, g.Config().Control.AllowUser) {
			this.SetSession("loginFailed", "2")
			logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: user %s is not allowed", clientIP, logtime, username))
			this.Ctx.Redirect(302, fmt.Sprintf("/auth/login?target=%s", target))
			return
		}
	}

	if len(g.Config().Control.Users) > 0 {
		user, err := g.Config().Control.GetUser(username)
		if err == nil {
			lastTry := user.LastTry
			user.LastTry = time.Now()
			if user.Num > 0 && user.Num > user.TryNum && time.Now().Unix()-lastTry.Unix() < int64(user.TryNum*60) {
				this.SetSession("loginFailed", "4")
				logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: %s", clientIP, logtime, "尝试次数过多"))
				this.Ctx.Redirect(302, fmt.Sprintf("/auth/login?target=%s", target))
				return
			}

			if user.Num > 0 && time.Now().Unix()-lastTry.Unix() < int64(user.Num*10) {
				this.SetSession("loginFailed", "5")
				logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: %s", clientIP, logtime, "尝试过于频繁"))
				this.Ctx.Redirect(302, fmt.Sprintf("/auth/login?target=%s", target))
				return
			}

			if user.PassWord == password {
				//登录成功设置session
				user.Num = 0
				if target == "" || target == "/auth/login" {
					logs.Warning(fmt.Sprintf("%s - - [%s] Login Failed: Missing X-Target", clientIP, logtime))
					this.Ctx.Redirect(302, "/")
				}
				this.SetSession("uname", username)
				logs.Notice(fmt.Sprintf("%s - %s [%s] Login Successed", clientIP, username, logtime))
				this.Ctx.Redirect(302, target)
			} else {
				user.Num = user.Num + 1
				this.SetSession("loginFailed", "1")
				logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: %s, num:%d", clientIP, logtime, "用户名密码错误", user.Num))
				this.Ctx.Redirect(302, fmt.Sprintf("/auth/login?target=%s", target))
			}
			return
		}
	}

	err := utils.LDAP_Auth(g.Config().Ldap, username, password)
	if err == nil {
		//登录成功设置session

		if target == "" || target == "/auth/login" {
			logs.Warning(fmt.Sprintf("%s - - [%s] Login Failed: Missing X-Target", clientIP, logtime))
			this.Ctx.Redirect(302, "/")
		}
		this.SetSession("uname", username)
		logs.Notice(fmt.Sprintf("%s - %s [%s] Login Successed", clientIP, username, logtime))
		this.Ctx.Redirect(302, target)
	} else {
		this.SetSession("loginFailed", "1")
		logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: %s", clientIP, logtime, err.Error()))
		this.Ctx.Redirect(302, fmt.Sprintf("/auth/login?target=%s", target))
	}
}
