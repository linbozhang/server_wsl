// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"wsl/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/rawbrand",&controllers.RawBrandController{})
	beego.Router("v1/vuforia",&controllers.VuforiaController{})
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/artarget",
			beego.NSInclude(
				&controllers.ArtargetController{},
			),
		),

		beego.NSNamespace("/brand",
			beego.NSInclude(
				&controllers.BrandController{},
			),
		),
		// beego.NSNamespace("/rawbrand",
		// 	beego.NSInclude(
		// 		&controllers.RawBrandController{},
		// 	),
		// ),
		beego.NSNamespace("/shwo",
			beego.NSInclude(
				&controllers.ShwoController{},
			),
		),

		beego.NSNamespace("/sort",
			beego.NSInclude(
				&controllers.SortController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
