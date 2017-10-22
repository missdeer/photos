package routers

import (
	"photos/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/img", "static/img")
	c := &controllers.MainController{}
	beego.Router("/", c)
	beego.Router("/p/:path", c, "get:GetPage")
	beego.Router("/i/:path", c, "get:GetImage")
	beego.Router("/s/:path", c, "get:GetSmallImage")
	beego.Router("/b/:path", c, "get:GetBigImage")
	beego.Router("/v/:path", c, "get:GetVideo")
}
