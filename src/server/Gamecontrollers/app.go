package Gamecontrollers

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)


// baseController represents base router for all other app routers.
// It implemented some methods for the same implementation;
// thus, it will be embedded into other routers.
type baseController struct {
	beego.Controller // Embed struct that has stub implementation of the interface.
	i18n.Locale      // For i18n usage when process data and render template.
}

// Prepare implemented Prepare() method for baseController.
// It's used for language option check and setting.
func (this *baseController) Prepare() {
	// Reset language option.
	this.Lang = "" // This field is from i18n.Locale.

	// 1. Get language information from 'Accept-Language'.
	al := this.Ctx.Request.Header.Get("Accept-Language")
	if len(al) > 4 {
		al = al[:5] // Only compare first 5 letters.
		if i18n.IsExist(al) {
			this.Lang = al
		}
	}

	// 2. Default language is English.
	if len(this.Lang) == 0 {
		this.Lang = "en-US"
	}

	// Set template level language option.
	this.Data["Lang"] = this.Lang
}

// AppController handles the welcome screen that allows user to pick a technology and username.
type AppController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for AppController.
func (this *AppController) Get() {
	this.TplName = "welcome.html"
}

// Join method handles POST requests for AppController.
func (this *AppController) Join() {
	// Get form value.
	uname := this.GetString("uname")

	// Check valid.
	if len(uname) == 0 {
		this.Redirect("/IM/", 302)
		return
	}


	this.Redirect("/IM/ws?uname="+uname, 302)


	// Usually put return after redirect.
	return
}
