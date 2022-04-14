package repo

import (
	"sale/internal/app/client"
)

type OfficeRepo interface {
	GetOfficesUUIDsByRegion(regionID uint) []string
	GetMainStoreOfficeInRegion(regionID uint) string
}

type RORepo struct { // мс регионов и офисов РиО
	client client.OfficerClient
}

func (r *RORepo) GetOfficesUUIDsByRegion(regionID uint) []string {
	return r.client.GetOfficesUUIDsByRegion(regionID)
}

func (r *RORepo) GetMainStoreOfficeInRegion(regionID uint) string {
	return r.client.GetMainStoreOfficeInRegion(regionID)
}
