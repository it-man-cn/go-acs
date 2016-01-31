package routers

import (
	"github.com/astaxie/beego"
	"go-acs/acs/controllers"
)

func init() {
	beego.Router("/tr069", &controllers.MainController{})
}
