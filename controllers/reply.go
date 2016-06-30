package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type ReplyController struct {
	beego.Controller
}

func (this *ReplyController) Add() {
	name := this.Input().Get("nickname")
	content := this.Input().Get("content")
	id := this.Input().Get("tid")

	err := models.AddReply(id, name, content)
	if err != nil {
		this.Redirect("/", 302)
		return
	}

	this.Redirect("/topic", 302)
	return
}

func (this *ReplyController) Del() {
	tid := this.Input().Get("tid")
	rid := this.Input().Get("rid")

	err := models.DelReply(tid, rid)
	if err != nil {
		this.Redirect("/", 302)
		return
	}

	this.Redirect("/topic", 302)
	return
}
