package controllers

import "github.com/astaxie/beego"


type MyContoller struct {
	beego.Controller
}

//这个函数用来接收sdk
func (controller *MyContoller) Init() {

}

func (controller *MyContoller) Index() {
	controller.TplName = "index.tpl"
}

func (controller *MyContoller) store() {
	_ = controller.Ctx.Input.RequestBody

}

