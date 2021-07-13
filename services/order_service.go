package services

import (
	"winson-product/datamodels"
	"winson-product/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order)error
	InsertOrder(*datamodels.Order) (int64,error)
	GetAllOrder()([]*datamodels.Order,error)
	GetAllOrderInfo()(map[int]map[string]string ,error)
	InsertOrderByMessage(*datamodels.Message) (int64, error)
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService{
	return &OrderService{OrderRepository: repository}
}

func (o *OrderService) GetOrderByID(orderID int64) (order *datamodels.Order, err error) {
	return o.OrderRepository.SelectByKey(orderID)
}

func (o *OrderService) DeleteOrderByID(orderID int64) (isOk bool) {
	return o.OrderRepository.Delete(orderID)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) (err error) {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (orderID int64, err error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder()(orderArray []*datamodels.Order, err error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo()(infoMap map[int]map[string]string ,err error) {
	return o.OrderRepository.SelectAllWithInfo()
}

//根据消息创建订单
func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error){
	order := &datamodels.Order{
		UserId: message.UserID,
		ProductId: message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}