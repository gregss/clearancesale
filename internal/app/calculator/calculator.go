package calculator

import (
	"errors"
	"math"

	"github.com/rs/zerolog/log"
	"sale/internal/app/repo"
	"sale/internal/models"
)

const (
	maxRegionSalePriceRate = 0.99 // не более 0.99 от продажной цены в регионе
	minPurchasePriceRate   = 0.5  // не менее половины закупочной цены
)

var ErrNullPrice = errors.New("нулевая цена")
var ErrNullSalePriceError = errors.New("нулевая продажная цена в регионе")
var ErrMaxLimitSalePriceError = errors.New("цена выше чем максимально возможная")
var ErrNullPurchasePriceError = errors.New("нулевая закупочная цена")
var ErrMinSalePriceError = errors.New("цена ниже чем минимально возможная")

// todo возможно это должна быть отдельная табличка в бд и вынести в другое место.

type CalcValues struct {
	StockAvailable uint
	Available      uint
	Status         uint8
}

type ParamsSaver interface {
	SaveParams(uint32, *CalcValues) error
}

type Calculator struct {
	storage         ParamsSaver
	stockCalculator StockCalculator
	limitApplier    LimitApplier
	priceRepo       repo.PriceGetter
}

func NewCalculator(product *models.SaleProduct, storage ParamsSaver) (*Calculator, error) {
	stockAvailableCalculator, err := NewStockAvailableCalculator(product.SaleID)
	if err != nil {
		return nil, err
	}

	return &Calculator{
		storage:         storage,
		stockCalculator: stockAvailableCalculator,
		limitApplier:    &LimitCalculator{storage: storage.(PurchasedCountGetterInterface)},
		priceRepo:       &repo.PriceRepo{},
	}, nil
}

func (c *Calculator) Calculate(product *models.SaleProduct) (*CalcValues, error) {
	// выход из распродажи необратим
	if product.Status == models.SaleProductIsInactive {
		return nil, nil
	}

	if err := c.checkPrice(product); err != nil {
		log.Warn().Err(err)
		return &CalcValues{}, nil
	}

	stockAvailable := c.stockCalculator.Calculate(product.NomenclatureUUID, product.RegionID)
	available := c.limitApplier.ApplyLimits(product, stockAvailable)
	var status uint8
	if available > 0 {
		status = 1
	}

	return &CalcValues{
		StockAvailable: stockAvailable,
		Available:      available,
		Status:         status,
	}, nil
}

func (c *Calculator) checkPrice(product *models.SaleProduct) error {
	if product.Price == 0 {
		return ErrNullPrice
	}

	maxRegionSalePrice := c.getMaxRegionSalePrice(product.RegionID, product.NomenclatureUUID)
	if maxRegionSalePrice == 0 {
		return ErrNullSalePriceError
	}
	if product.Price > maxRegionSalePrice {
		return ErrMaxLimitSalePriceError
	}

	minPurchasePrice := c.getMinPurchasePriceRate(product.NomenclatureUUID)
	if minPurchasePrice == 0 {
		return ErrNullPurchasePriceError
	}
	if product.Price < minPurchasePrice {
		return ErrMinSalePriceError
	}

	return nil
}

func (c *Calculator) getMaxRegionSalePrice(regionID uint, nomenclUUID string) uint {
	return uint(math.Round(float64(c.priceRepo.GetSalePrice(regionID, nomenclUUID)) * maxRegionSalePriceRate))
}

func (c *Calculator) getMinPurchasePriceRate(nomenclUUID string) uint {
	return uint(math.Round(float64(c.priceRepo.GetPurchasePrice(nomenclUUID)) * minPurchasePriceRate))
}
