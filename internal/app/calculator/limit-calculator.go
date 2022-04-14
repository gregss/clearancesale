package calculator

import (
	"github.com/rs/zerolog/log"
	"sale/internal/models"
	"sale/internal/pkg/math"
)

type LimitApplier interface {
	ApplyLimits(product *models.SaleProduct, stockAvailable uint) uint
}

type PurchasedCountGetterInterface interface {
	GetPurchasedCount(product *models.SaleProduct) (uint, error)
}

type LimitCalculator struct {
	storage PurchasedCountGetterInterface
}

func (c *LimitCalculator) ApplyLimits(product *models.SaleProduct, stockAvailable uint) uint {
	available := stockAvailable
	if product.MaxCount > 0 {
		maxCountLimitAvailable := c.maxCountLimit(product)
		// todo log.info причина выхода из распродажи
		available = math.Min(available, maxCountLimitAvailable)
	}

	if product.MaxOrderCount > 0 {
		// todo log.info причина выхода из распродажи
		available = math.Min(available, product.MaxOrderCount)
	}

	return available
}

func (c *LimitCalculator) maxCountLimit(product *models.SaleProduct) uint {
	purchasedCount, err := c.storage.GetPurchasedCount(product)
	if err != nil {
		log.Err(err).Msg("error get purchased quantity")
		return 0
	}

	if purchasedCount > product.MaxCount {
		return 0
	}

	return product.MaxCount - purchasedCount
}
