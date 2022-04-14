package storage

import (
	"sale/internal/app/calculator"
	"sale/internal/models"
)

type Storage interface {
	CreateSale(*models.Sale) error
	AddProduct(*models.SaleProduct) error
	EditSale(*models.Sale) error
	EditProduct(*models.SaleProduct) error
	SaveParams(uint32, *calculator.CalcValues) error
	AddSell(*models.Sell) error
	GetSaleProduct(nomenclUUID string, regionID uint) (*models.SaleProduct, error)
}
