package main

import (
	"fmt"
	"winson-product/common"
	"winson-product/rabbitmq"
	"winson-product/repositories"
	"winson-product/services"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	//创建product数据库操作实例
	product := repositories.NewProductManager("product",db)
	//创建product service
	productService := services.NewProductService(product)
	//创建order数据库实例
	order := repositories.NewOrderManagerRepository("order",db)
	//创建order service
	orderService := services.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("winsonProduct")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
