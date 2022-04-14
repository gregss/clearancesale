package repo

import (
	"sale/internal/app/client"
)

type StockRepo interface {
	GetRestsByOffices(officesUUIDs []string, nomenclUUID string) uint
}

type ShcatRepo struct {
	client client.ShcatClient
}

func (r *ShcatRepo) GetRestsByOffices(officesUUIDs []string, nomenclUUID string) uint {
	return r.client.GetRestsByOffices(officesUUIDs, nomenclUUID)
}
