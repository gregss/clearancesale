package repo

import (
	"sale/internal/app/client"
)

type PriceGetter interface {
	GetSalePrice(regionID uint, nomenclUUID string) uint
	GetPurchasePrice(nomenclUUID string) uint
}

type PriceRepo struct {
	client client.ScroogeClient
}

func (r *PriceRepo) GetSalePrice(regionID uint, nomenclUUID string) uint {
	return r.client.GetSalePrice(regionID, nomenclUUID)
}

func (r *PriceRepo) GetPurchasePrice(nomenclUUID string) uint {
	return r.client.GetPurchasePrice(nomenclUUID)
}
