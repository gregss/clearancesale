package models

import (
	"time"
)

// todo у констант область видимости для всех моделей.
const (
	SaleProductIsActive   = 1
	SaleProductIsInactive = 0
)

type SaleProduct struct {
	ID               uint32
	UUID             string
	SaleID           uint
	NomenclatureUUID string
	RegionID         uint // todo main region
	MaxCount         uint // максимальное кол-во которое можно продать по распродаже
	MaxOrderCount    uint // максимальное кол-во которое можно продать в одном заказе
	Price            uint
	Type             uint8
	IsFeed           bool // Признак отображения цены в фидах'
	StockAvailable   uint // Доступное кол-во остатков
	Available        uint // Доступное кол-во в распродаже с учетом остатков и лимитов (главное вычисляемое поле)
	Status           uint8
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedBy        string
	UpdatedAt        time.Time
}
