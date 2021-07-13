package repositories

import (
	"database/sql"
	"strconv"
	"winson-product/common"
	"winson-product/datamodels"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64,error)
	Delete(int64) bool
	Update(*datamodels.Order)error
	SelectByKey(int64)(*datamodels.Order,error)
	SelectAll()([]*datamodels.Order,error)
	SelectAllWithInfo()(map[int]map[string]string ,error)
}

type OrderManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository{
	return &OrderManagerRepository{table: table, mysqlConn: sql}
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Conn(); err != nil {
		return 0,err
	}

	sql := "insert `order` set userID=?,productID=?,orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return productID, err
	}

	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (o *OrderManagerRepository) Delete(orderID int64) (isOk bool) {
	if err := o.Conn(); err != nil {
		return
	}

	sql := "delete from "+ o.table +" where ID=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}

	_, err = stmt.Exec(orderID)
	if err != nil {
		return
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) (err error) {
	if err = o.Conn(); err != nil {
		return err
	}

	sql := "update "+ o.table +" set userID=?,productID=?,orderStatus=? where ID=?"+ strconv.FormatInt(order.ID, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return  err
	}

	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return err
	}
	return
}

func (o *OrderManagerRepository) SelectByKey(orderID int64)(order *datamodels.Order,err error){
	if err = o.Conn(); err != nil {
		return &datamodels.Order{},err
	}

	sql := "select * from "+ o.table +" where ID="+ strconv.FormatInt(orderID, 10)
	row, err := o.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.Order{}, err
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}

	order = &datamodels.Order{}
	common.DataToStructByTagSql(result,order)
	return
}

func (o *OrderManagerRepository) SelectAll()(orderArray []*datamodels.Order, err error){
	if err = o.Conn(); err != nil {
		return nil,err
	}

	sql := "select * from "+ o.table
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	for _, v := range result{
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo()(OrderMap map[int]map[string]string , err error){
	if err = o.Conn(); err != nil {
		return nil,err
	}

	sql := "select o.ID, p.productName, o.orderStatus from winson.order as o left join winson.product as p on o.productID=p.ID"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	return common.GetResultRows(rows), err
}

