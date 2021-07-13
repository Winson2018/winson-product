package main

import (
	"fmt"
	"strconv"
	"winson-product/common"
	"winson-product/datamodels"
)

func testmain() {
	data := map[string]string{"ID":"1","productName":"winson测试结构体",
		"productNum":"2","productImage":"123","productUrl":"http://url"}

	product := &datamodels.Product{}
	common.DataToStructByTagSql(data,product)

	fmt.Println("id="+ strconv.FormatInt(product.ID, 10))
	fmt.Println("productName="+ product.ProductName)
	fmt.Println("productNum="+ strconv.FormatInt(product.ProductNum,10))
	fmt.Println("productImage="+ product.ProductImage)
	fmt.Println("productUrl="+ product.ProductUrl)
}
