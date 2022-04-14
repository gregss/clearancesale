package models

import (
	"time"
)

// todo возможно это нужно вынести отсюда.
const (
	InRegionStopFactorID                    = 1
	WithoutMainDistributeCenterStopFactorID = 2
	MainDistributeCenterStopFactorID        = 3
	LogisticChainStopFactorID               = 4
)

// State.
const (
	SaleIsInactive = 0 // zero value
	SaleIsActive   = 1
)

// Contractor type.
const (
	All = 0 // zero value
	B2b = 1
	B2c = 2
)

const (
	BlackFriday       = 1 // черная пятница
	ProductOfTheMonth = 2 // товар месяца
	ShopSale          = 3 // распродажа в магазине
	SuperPpriceSale   = 4 // супер цена
)

type Sale struct {
	ID             uint   // id vs uuid
	UUID           string // пока нужно для обменов
	Name           string // `json:"name"`
	StartDate      time.Time
	EndDate        time.Time
	StopFactor     uint8 // todo отдельный тип?
	ContractorType uint8
	Type           uint8
	State          uint8
	CreatedBy      string // ldap uuid
	CreatedAt      time.Time
}
