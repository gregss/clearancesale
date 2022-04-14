package client

import (
	"net/http"
)

// todo заменить на sdk мс.
const (
	getOfficesUUIDsByRegion = "/api/v1/json/getOfficesUUIDsByRegion"
)

type OfficerClient struct {
	client http.Client
}

// GetOfficesUUIDsByRegion запрос к стороннему мс офисов на получение uuids офисов в регионе.
func (c *OfficerClient) GetOfficesUUIDsByRegion(regionID uint) []string {
	_ = regionID
	_, _ = c.client.Get(getOfficesUUIDsByRegion)
	return []string{"office1"}
}

func (c *OfficerClient) GetMainStoreOfficeInRegion(regionID uint) string {
	_ = regionID
	_, _ = c.client.Get(getOfficesUUIDsByRegion)
	return "office2"
}
