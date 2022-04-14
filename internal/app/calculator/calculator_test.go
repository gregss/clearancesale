package calculator

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"sale/internal/models"
	mock_repo "sale/tests/mock"
)

type paramsSaver struct{}

func (s *paramsSaver) SaveParams(id uint32, v *CalcValues) error {
	return nil
}

type stockCalculator struct{}

func (c *stockCalculator) Calculate(nomenclUUID string, regionID uint) uint {
	return 6
}

type limitApplier struct{}

func (a *limitApplier) ApplyLimits(product *models.SaleProduct, stockAvailable uint) uint {
	return 6
}

/*type priceGetter struct{}

func (g *priceGetter) GetSalePrice(regionID uint, nomenclUUID string) uint {
	return 100
}
func (g *priceGetter) GetPurchasePrice(nomenclUUID string) uint {
	return 50
}*/

func TestCalculator(t *testing.T) {
	t.Skip()
	contrl := gomock.NewController(t)
	defer contrl.Finish()
	priceGetterMock := mock_repo.NewMockPriceGetter(contrl)

	calculator := &Calculator{
		storage:         &paramsSaver{},
		stockCalculator: &stockCalculator{},
		limitApplier:    &limitApplier{},
		priceRepo:       priceGetterMock,
	}

	_, err := calculator.Calculate(&models.SaleProduct{
		Price:  1,
		Status: models.SaleProductIsInactive,
	})

	require.NoError(t, err)

	priceGetterMock.
		EXPECT().
		GetSalePrice(gomock.All(), gomock.All()).
		Return(uint(0))

	_, err = calculator.Calculate(&models.SaleProduct{
		Price:  70,
		Status: models.SaleProductIsActive,
	})

	require.ErrorIs(t, err, ErrNullSalePriceError)

	priceGetterMock.
		EXPECT().
		GetSalePrice(gomock.All(), gomock.All()).
		Return(uint(100))

	priceGetterMock.
		EXPECT().
		GetPurchasePrice(gomock.All()).
		Return(uint(0))

	_, err = calculator.Calculate(&models.SaleProduct{
		Price:  70,
		Status: models.SaleProductIsActive,
	})

	require.ErrorIs(t, err, ErrNullPurchasePriceError)

	tests := []struct {
		tcase string
		price uint
		err   error
	}{
		{
			tcase: "in limit price",
			price: 70,
			err:   nil,
		},
		{
			tcase: "over max price",
			price: 200,
			err:   ErrMaxLimitSalePriceError,
		},
		{
			tcase: "over min price",
			price: 10,
			err:   ErrMinSalePriceError,
		},
	}

	priceGetterMock.
		EXPECT().
		GetSalePrice(gomock.All(), gomock.All()).
		Return(uint(100)).AnyTimes()

	priceGetterMock.
		EXPECT().
		GetPurchasePrice(gomock.All()).
		Return(uint(50)).AnyTimes()

	for _, test := range tests {
		test := test
		t.Run(test.tcase, func(t *testing.T) {
			_, err := calculator.Calculate(&models.SaleProduct{
				Price:  test.price,
				Status: models.SaleProductIsActive,
			})

			require.ErrorIs(t, err, test.err)
		})
	}
}
