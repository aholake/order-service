package db

import (
	"fmt"

	"github.com/aholake/order-service/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(connectionUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(connectionUrl), &gorm.Config{})

	if openErr != nil {
		return nil, fmt.Errorf("DB connection error %v", openErr)
	}
	error := db.AutoMigrate(Order{}, OrderItem{})
	if error != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a Adapter) Get(id string) (domain.Order, error) {
	var orderEntity Order
	tx := a.db.First(&orderEntity, id)
	var orderItems []domain.OrderItem
	for _, oi := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: oi.ProductCode,
			UnitPrice:   oi.UnitPrice,
			Quantity:    oi.Quantity,
		})
	}
	return domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.Unix(),
	}, tx.Error
}

func (a Adapter) Save(order *domain.Order) error {
	var orderItemEntities []OrderItem
	for _, oi := range order.OrderItems {
		orderItemEntities = append(orderItemEntities, OrderItem{
			ProductCode: oi.ProductCode,
			Quantity:    oi.Quantity,
			UnitPrice:   oi.UnitPrice,
		})
	}

	orderModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItemEntities,
	}

	tx := a.db.Create(&orderModel)
	if tx.Error == nil {
		order.ID = int64(orderModel.ID)
	}
	return tx.Error
}
