package calculator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type stringMock struct{}

// todo заменить на golang/mock или подобный.
func (c *stringMock) GetOfficesUUIDsByRegion(regionID uint) []string {
	return []string{"officeUUID1", "officeUUID2"}
}

func (c *stringMock) GetMainStoreOfficeInRegion(regionID uint) string {
	return "officeUUID3"
}

type stockRepoMock struct{}

func (c *stockRepoMock) GetRestsByOffices(officesUUIDs []string, nomenclUUID string) uint {
	stockData := map[string]map[string]uint{
		"officeUUID1": {"nomenclUUID1": 1},
		"officeUUID2": {"nomenclUUID1": 2, "nomenclUUID2": 1},
		"officeUUID3": {"nomenclUUID1": 1},
		"officeUUID4": {"nomenclUUID1": 1},
		"officeUUID5": {"nomenclUUID1": 1},
	}

	var rest uint
	for _, officeUUID := range officesUUIDs {
		if r, ok := stockData[officeUUID][nomenclUUID]; ok {
			rest += r
		}
	}

	return rest
}

type logisticRepoMock struct{}

func (c *logisticRepoMock) GetLogisticOfficesTo(officesToUUID string) []string {
	logisticData := map[string][]string{
		"officeUUID3": {"officeUUID4", "officeUUID5"},
		"officeUUID4": {"officeUUID5"},
	}

	if r, ok := logisticData[officesToUUID]; ok {
		return r
	}

	return nil
}

type providerRepoMock struct{}

func (c *providerRepoMock) HasProvider(nomenclUUID string) bool {
	providerData := map[string]bool{
		"nomenclUUID1": true,
		"nomenclUUID2": false,
	}

	if r, ok := providerData[nomenclUUID]; ok {
		return r
	}

	return false
}

func TestInRegionStopFactorCalculator(t *testing.T) {
	calculator := &InRegionStopFactorCalculator{string: &stringMock{}, stockRepo: &stockRepoMock{}}

	require.Equal(t, uint(3), calculator.Calculate("nomenclUUID1", 1))
	require.Equal(t, uint(1), calculator.Calculate("nomenclUUID2", 1))
}

func TestMainDistributeCenterStopFactorCalculator(t *testing.T) {
	calculator := &MainDistributeCenterStopFactorCalculator{
		string:       &stringMock{},
		stockRepo:    &stockRepoMock{},
		logisticRepo: &logisticRepoMock{},
	}

	require.Equal(t, uint(5), calculator.Calculate("nomenclUUID1", 1))
	require.Equal(t, uint(1), calculator.Calculate("nomenclUUID2", 1))
}

func TestLogisticChainStopFactorCalculator(t *testing.T) {
	calculator := &LogisticChainStopFactorCalculator{
		providerRepo: &providerRepoMock{},
		sfcalc: &MainDistributeCenterStopFactorCalculator{
			string:       &stringMock{},
			stockRepo:    &stockRepoMock{},
			logisticRepo: &logisticRepoMock{},
		},
	}

	require.Equal(t, uint(9999999), calculator.Calculate("nomenclUUID1", 1))
	require.Equal(t, uint(1), calculator.Calculate("nomenclUUID2", 1))
}
