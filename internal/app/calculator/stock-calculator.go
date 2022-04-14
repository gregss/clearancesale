package calculator

import (
	"errors"

	"sale/internal/app/repo"
	"sale/internal/models"
)

type StockCalculator interface {
	Calculate(nomenclUUID string, regionID uint) uint
}

type StockAvailableCalculator struct {
	stopFactorCalculator StopFactorCalculator
}

func (c *StockAvailableCalculator) Calculate(nomenclUUID string, regionID uint) uint {
	return c.stopFactorCalculator.Calculate(nomenclUUID, regionID)
}

func NewStockAvailableCalculator(stopFactor uint) (*StockAvailableCalculator, error) {
	var stopFactorCalculator StopFactorCalculator

	switch stopFactor {
	case models.InRegionStopFactorID:
		stopFactorCalculator = &InRegionStopFactorCalculator{
			string:    &repo.RORepo{},
			stockRepo: &repo.ShcatRepo{},
		}
	case models.WithoutMainDistributeCenterStopFactorID:
		stopFactorCalculator = &WithoutMainDistributeCenterStopFactorCalculator{
			string:       &repo.RORepo{},
			stockRepo:    &repo.ShcatRepo{},
			logisticRepo: &repo.SusaninRepo{},
		}
	case models.MainDistributeCenterStopFactorID:
		stopFactorCalculator = &MainDistributeCenterStopFactorCalculator{
			string:       &repo.RORepo{},
			stockRepo:    &repo.ShcatRepo{},
			logisticRepo: &repo.SusaninRepo{},
		}
	case models.LogisticChainStopFactorID:
		stopFactorCalculator = &LogisticChainStopFactorCalculator{
			providerRepo: &repo.PDMRepo{},
			sfcalc: &MainDistributeCenterStopFactorCalculator{
				string:       &repo.RORepo{},
				stockRepo:    &repo.ShcatRepo{},
				logisticRepo: &repo.SusaninRepo{},
			},
		}
	default:
		return nil, errors.New("неизвестный стопфактор")
	}

	return &StockAvailableCalculator{stopFactorCalculator: stopFactorCalculator}, nil
}

// todo возможно в отедльный файл в том же namespace.
const unlimit = 9999999

type StopFactorCalculator interface {
	Calculate(nomenclUUID string, regionID uint) uint
}

type InRegionStopFactorCalculator struct {
	string    repo.OfficeRepo
	stockRepo repo.StockRepo
}
type WithoutMainDistributeCenterStopFactorCalculator struct {
	string       repo.OfficeRepo
	stockRepo    repo.StockRepo
	logisticRepo repo.LogisticRepo
}
type MainDistributeCenterStopFactorCalculator struct {
	string       repo.OfficeRepo
	stockRepo    repo.StockRepo
	logisticRepo repo.LogisticRepo
}
type LogisticChainStopFactorCalculator struct {
	providerRepo repo.ProviderRepo
	sfcalc       StopFactorCalculator
}

func (c *InRegionStopFactorCalculator) Calculate(nomenclUUID string, regionID uint) uint {
	officesUUIDs := c.string.GetOfficesUUIDsByRegion(regionID)

	if len(officesUUIDs) == 0 {
		return 0
	}

	return c.stockRepo.GetRestsByOffices(officesUUIDs, nomenclUUID)
}

func (c *WithoutMainDistributeCenterStopFactorCalculator) Calculate(nomenclUUID string, regionID uint) uint {
	// todo
	return 0
}

func (c *MainDistributeCenterStopFactorCalculator) Calculate(nomenclUUID string, regionID uint) uint {
	officesUUIDs := append(
		c.logisticRepo.GetLogisticOfficesTo(c.string.GetMainStoreOfficeInRegion(regionID)),
		c.string.GetOfficesUUIDsByRegion(regionID)...,
	)

	if len(officesUUIDs) == 0 {
		return 0
	}

	return c.stockRepo.GetRestsByOffices(officesUUIDs, nomenclUUID)
}

func (c *LogisticChainStopFactorCalculator) Calculate(nomenclUUID string, regionID uint) uint {
	if c.providerRepo.HasProvider(nomenclUUID) {
		return unlimit
	}

	// если поставщика нет, считаем остатки по лог. цепочке
	return c.sfcalc.Calculate(nomenclUUID, regionID)
}
