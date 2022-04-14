package repo

import (
	"sale/internal/app/client"
)

type LogisticRepo interface {
	GetLogisticOfficesTo(officesToUUID string) []string
}

type SusaninRepo struct {
	client client.SusaninClient
}

func (r *SusaninRepo) GetLogisticOfficesTo(officesToUUID string) []string {
	return r.client.GetLogisticOfficesTo(officesToUUID)
}
