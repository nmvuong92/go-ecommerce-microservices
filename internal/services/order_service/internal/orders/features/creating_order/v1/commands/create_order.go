package createOrderCommandV1

import (
	"time"

	uuid "github.com/satori/go.uuid"

	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator
type CreateOrder struct {
	OrderId         uuid.UUID             `validate:"required"`
	ShopItems       []*dtosV1.ShopItemDto `validate:"required"`
	AccountEmail    string                `validate:"required,email"`
	DeliveryAddress string                `validate:"required"`
	DeliveryTime    time.Time             `validate:"required"`
	CreatedAt       time.Time             `validate:"required"`
}

func NewCreateOrder(
	shopItems []*dtosV1.ShopItemDto,
	accountEmail, deliveryAddress string,
	deliveryTime time.Time,
) *CreateOrder {
	return &CreateOrder{
		OrderId:         uuid.NewV4(),
		ShopItems:       shopItems,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		DeliveryTime:    deliveryTime,
		CreatedAt:       time.Now(),
	}
}
