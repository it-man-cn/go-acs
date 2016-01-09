package routers

import (
	"github.com/astaxie/beego"
	"go-acs/controllers"
)

func init() {
	beego.Router("/tr069", &controllers.MainController{})
	beego.EnableAdmin = true
	beego.AdminHttpPort = 8888
}
