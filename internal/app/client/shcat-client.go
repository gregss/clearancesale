package client

import (
	"net/http"
)

// todo заменить на sdk мс.
const (
	getRestsByOffices = "/api/v1/json/getRestsByOffices"
)

type ShcatClient struct { // Shcat - мс остатков
	client http.Client
}

func (c *ShcatClient) GetRestsByOffices(officesUUIDs []string, nomenclUUID string) uint {
	_ = officesUUIDs
	_ = nomenclUUID
	_, _ = c.client.Get(getRestsByOffices)
	return 10
}
