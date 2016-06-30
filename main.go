package main

import (
	_ "beeblog/routers"
	"beeblog/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	// 开启 ORM 调试模式
	models.RegisterDB()
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)

	beego.Run()
}
