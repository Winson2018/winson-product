package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"winson-product/common"
	"winson-product/fronted/middleware"
	"winson-product/fronted/web/controllers"
	"winson-product/rabbitmq"
	"winson-product/repositories"
	"winson-product/services"
)

func main() {

	app := iris.New()
	app.Logger().SetLevel("debug")
	tmplate := iris.HTML("./fronted/web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmplate)
	//app.StaticWeb("/public", "./fronted/web/public")
	app.HandleDir("/public", "./fronted/web/public")
	//访问生成好的HTML静态文件
	app.HandleDir("/html", "./fronted/web/htmlProductShow")
	app.OnAnyErrorCode(func(ctx iris.Context){
		ctx.ViewData("message", ctx.Values().GetStringDefault("message","访问的页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {
		log.Println(err)
	}

/*	sess := sessions.New(sessions.Config{
		Cookie: "AdminCookie",
		Expires: 600 * time.Minute,
	})*/

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//注册user控制器
	user := repositories.NewUserRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	//userPro.Register(userService, ctx, sess.Start)
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	rabbitmq := rabbitmq.NewRabbitMQSimple("winsonProduct")

	//注册product控制器
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct)
	pro.Register(productService, orderService, rabbitmq)
	//pro.Register(productService, orderService)
	//pro.Register(productService, sess.Start)
	pro.Handle(new(controllers.ProductController))


	app.Run(
		iris.Addr("0.0.0.0:8082"),
		//iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		)
}
