package routers

import (
	"beeblog/controllers"
	"beeblog/models"
	"github.com/astaxie/beego"
)

func init() {
	// 注册数据库
	models.RegisterDB()

	beego.Router("/", &controllers.HomeController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/login", &controllers.LoginController{})
}
