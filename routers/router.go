package routers
//路由层
import (
	"github.com/astaxie/beego"
	"myfabric/controllers"
)

func init() {
	beego.Router("/", &controllers.MyContoller{}, "get,post:Index")

}
