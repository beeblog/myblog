package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (this *TopicController) Get() {
	this.Data["IsTopic"] = "true"
	this.TplName = "topic.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	var err error
	this.Data["Topics"], err = models.GetAllTopics()

	if err != nil {
		beego.Error(err)
	}
}

func (this *TopicController) Add() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.TplName = "topic_add.html"
}

func (this *TopicController) Post() {
	id := this.Input().Get("tid")
	title := this.Input().Get("title")
	content := this.Input().Get("content")
	category := this.Input().Get("category")

	if len(title) == 0 || len(content) == 0 {
		this.Redirect("/topic", 302)
		return
	}
	if len(id) == 0 {
		err := models.AddTopic(title, content, category)
		if err != nil {
			beego.Error(err)
			this.Redirect("/topic/add", 302)
			return
		}
	} else {
		err := models.ModifyTopic(id, title, content)
		if err != nil {
			beego.Error(err)
			this.Redirect("/topic/add", 302)
			return
		}
	}

	this.Redirect("/topic", 302)
}

func (this *TopicController) Modify() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	id := this.Input().Get("tid")

	this.TplName = "topic_modify.html"
	var err error
	topic, err := models.GetTopic(id)
	if err != nil {
		beego.Error(err)
		this.Redirect("/topic", 302)
		return
	}

	this.Data["Topic"] = topic
	this.Data["Tid"] = id
}

func (this *TopicController) View() {
	this.TplName = "topic_view.html"

	topic, err := models.GetTopic(this.Ctx.Input.Params()["0"])
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	this.Data["Topic"] = topic
}

func (this *TopicController) Delete() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	err := models.DeleteTopic(this.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}

	this.Redirect("/topic", 302)
}
