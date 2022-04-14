package repo

import (
	"sale/internal/app/client"
)

type ProviderRepo interface {
	HasProvider(nomenclUUID string) bool
}

type PDMRepo struct {
	client client.PdmClient
}

func (r *PDMRepo) HasProvider(nomenclUUID string) bool {
	return r.client.HasProvider(nomenclUUID)
}
