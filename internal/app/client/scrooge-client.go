package client

import (
	"net/http"
)

// todo заменить на sdk мс.
const (
	getMaxRegionSalePrice   = "/api/v1/json/getMaxRegionSalePrice"
	getMinPurchasePriceRate = "/api/v1/json/getMinPurchasePriceRate"
)

type ScroogeClient struct {
	client http.Client
}

func (c *ScroogeClient) GetSalePrice(regionID uint, nomenclUUID string) uint {
	_ = regionID
	_ = nomenclUUID
	_, _ = c.client.Get(getMaxRegionSalePrice)
	return 100
}

func (c *ScroogeClient) GetPurchasePrice(nomenclUUID string) uint {
	_ = nomenclUUID
	_, _ = c.client.Get(getMinPurchasePriceRate)
	return 50
}
