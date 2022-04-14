package client

import (
	"net/http"
)

// todo заменить на sdk мс.
const (
	getProvider = "/api/v1/json/getProvider"
)

type PdmClient struct { // PDM - мс номенклатур
	client http.Client
}

func (c *PdmClient) HasProvider(nomenclUUID string) bool {
	_ = nomenclUUID
	_, _ = c.client.Get(getProvider)
	return true
}
