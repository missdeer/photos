package controllers

import (
	"github.com/astaxie/beego"
)

// Photo struct represents a pohto
type Photo struct {
	Origin string
	Small  string
	Title  string
}

// MainController main controller
type MainController struct {
	beego.Controller
}

// Get return main page
func (c *MainController) Get() {
	var photos []Photo
	c.Data["Photos"] = photos
	c.Data["Title"] = "Our Photos"
	c.TplName = "index.tpl"
}

// GetImage return image
func (c *MainController) GetImage() {
	c.Ctx.Input.Param(":path")
}

// GetPage return a specified page
func (c *MainController) GetPage() {
	path := c.Ctx.Input.Param(":path")
	var photos []Photo
	c.Data["Photos"] = photos
	c.Data["Title"] = path
	c.TplName = "index.tpl"
}
