package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" winson:"ID"`
	ProductName  string `json:"ProductName" sql:"productName" winson:"ProductName"`
	ProductNum   int64 `json:"ProductNum" sql:"productNum" winson:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" winson:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" winson:"ProductUrl"`
}
