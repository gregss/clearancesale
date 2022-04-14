package client

import (
	"net/http"
)

// todo заменить на sdk мс.
const (
	getLogisticOfficesTo = "/api/v1/json/getLogisticOfficesTo"
)

type SusaninClient struct { // Shcat - мс остатков
	client http.Client
}

func (c *SusaninClient) GetLogisticOfficesTo(officesToUUID string) []string {
	_ = officesToUUID
	_, _ = c.client.Get(getLogisticOfficesTo)
	return []string{"office2", "office3"}
}
