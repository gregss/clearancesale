package calculator

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"sale/internal/models"
)

type PurchasedCountGetterMock struct{}

func (g *PurchasedCountGetterMock) GetPurchasedCount(product *models.SaleProduct) (uint, error) {
	return map[uint32]uint{
		1: 0,
		2: 2,
		3: 3,
	}[product.ID], nil
}

func TestLimitCalculator(t *testing.T) {
	calculator := &LimitCalculator{
		storage: &PurchasedCountGetterMock{},
	}

	product := &models.SaleProduct{
		ID:               1,
		NomenclatureUUID: uuid.NewString(),
		SaleID:           1,
		RegionID:         1,
		Price:            1,
		Type:             1,
		IsFeed:           true,
		StockAvailable:   1,
		Available:        1,
		Status:           1,
	}

	product.MaxCount = 5
	product.MaxOrderCount = 5

	// todo табличный тест
	require.Equal(t, uint(5), calculator.ApplyLimits(product, 10)) // достаточно на остатках
	require.Equal(t, uint(0), calculator.ApplyLimits(product, 0))  // нет на остатках
	require.Equal(t, uint(3), calculator.ApplyLimits(product, 3))  // остатков меньше лимитов

	product.MaxCount = 15
	product.MaxOrderCount = 5

	require.Equal(t, uint(5), calculator.ApplyLimits(product, 10))

	product.MaxCount = 5
	product.MaxOrderCount = 15

	require.Equal(t, uint(5), calculator.ApplyLimits(product, 10))
}
