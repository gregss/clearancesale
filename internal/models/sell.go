package models

import (
	"time"
)

// Sell данные использования распродажного товара в заказе.
type Sell struct {
	SaleProductID uint
	OrderUUID     string // товарная часть заказа
	Quantity      uint   // сколько использовано
	CreatedAt     time.Time
}
