//go:build integrations
// +build integrations

package tests

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"os"
	salev1 "sale/gen/api/public/v1"
	saleuiv1 "sale/gen/api/ui/v1"
	"sale/internal/models"
	"sale/internal/storage"
	"testing"
	"time"
)

type SaleSuite struct {
	suite.Suite
	uiClient     saleuiv1.SaleService
	publicClient salev1.SellingService
	storage      storage.Storage
}

func (s *SaleSuite) SetupSuite() {
	s.uiClient = saleuiv1.NewSaleServiceProtobufClient(
		"http://localhost:8081",
		&http.Client{},
	)

	s.publicClient = salev1.NewSellingServiceProtobufClient(
		"http://localhost:8080",
		&http.Client{},
	)

	//os.Getenv
}

func (s *SaleSuite) TearDownAllSuite() {

}

func (s *SaleSuite) SetupTest() {

}

func (s *SaleSuite) TestCreateSale() {
	s.createSale()
	s.createSaleProduct()
	time.Sleep(time.Second * 2)
	_, _ = s.sell()
	_, err := s.publicClient.Rest(context.Background(), &salev1.RestRequest{
		Rest: []*salev1.RestRequest_Rest{
			{
				NomenclatureUuid: "fd349237-e9b4-4aa0-acf3-101ef6b48000",
				RegionId:         1,
				OrderUuid:        "fd349237-e9b4-4aa0-acf3-101ef6b48555",
			},
			{
				NomenclatureUuid: "fd349237-e9b4-4aa0-acf3-101ef6b48111",
				RegionId:         1,
				OrderUuid:        "fd349237-e9b4-4aa0-acf3-101ef6b48555",
			},
		},
	})

	s.Require().NoError(err)
	// todo мок сервер внешних мс
	/*s.Require().Equal(salev1.RestResponse{
		Rest: []*salev1.RestResponse_Rest{
			{
				NomenclatureUuid: "fd349237-e9b4-4aa0-acf3-101ef6b48000",
				Sale: &salev1.RestResponse_Rest_Sale{
					SaleId:   1,
					SaleName: "Тестовая распродажа 1",
					Quantity: 1,
					Price:    100,
				},
			},
			{
				NomenclatureUuid: "fd349237-e9b4-4aa0-acf3-101ef6b48111",
				Sale: &salev1.RestResponse_Rest_Sale{
					SaleId:   1,
					SaleName: "Тестовая распродажа 1",
					Quantity: 1,
					Price:    200,
				},
			},
		},
	}, r)*/
}

func (s *SaleSuite) createSale() {
	startDate, _ := time.Parse("2006-01-02T15:04:05", "2022-05-01T00:00:00")
	endDate, _ := time.Parse("2006-01-02T15:04:05", "2022-05-30T00:00:00")

	_, err := s.uiClient.Create(context.Background(), &saleuiv1.CreateRequest{
		Name:           "Тестовая распродажа 1",
		StartDate:      timestamppb.New(startDate),
		EndDate:        timestamppb.New(endDate),
		StopFactor:     models.LogisticChainStopFactorID,
		ContractorType: models.All,
		Type:           models.BlackFriday,
	})
	if err != nil {
		fmt.Printf("error create sale: %v", err)
		os.Exit(1)
	}
	fmt.Println("sale created")
}

func (s *SaleSuite) createSaleProduct() {
	_, err := s.uiClient.CreateProducts(context.Background(), &saleuiv1.CreateProductsRequest{
		SaleId:   1,
		RegionId: []uint32{1, 2},
		SaleProducts: []*saleuiv1.CreateProductsRequest_SaleProducts{
			{
				NomenclatureUuid: "fd349237-e9b4-4aa0-acf3-101ef6b48000",
				Price:            1,
				MaxCount:         10,
				MaxOrderCount:    5,
				Type:             1,
				IsFeed:           true,
			},
			{
				NomenclatureUuid: "fd349237-e9b4-4aa0-acf3-101ef6b48111",
				Price:            200,
				MaxCount:         5,
				MaxOrderCount:    2,
				Type:             0,
				IsFeed:           true,
			},
		},
	})
	if err != nil {
		fmt.Printf("error create product sale: %v", err)
		os.Exit(1)
	}
	fmt.Println("sale created")
}

func (s *SaleSuite) sell() (*emptypb.Empty, error) {
	return s.publicClient.Sell(context.Background(), &salev1.SellRequest{
		Sell: []*salev1.SellRequest_Sell{
			{
				NomenclatureUuid: "fd349237-e9b4-4aa0-acf3-101ef6b48000",
				Quantity:         1,
				OrderUuid:        "fd349237-e9b4-4aa0-acf3-101ef6b48555",
			},
			{
				NomenclatureUuid: "fd349237-e9b4-4aa0-acf3-101ef6b48111",
				Quantity:         2,
				OrderUuid:        "fd349237-e9b4-4aa0-acf3-101ef6b48777",
			},
		},
	})
}

func TestSaleSuite(t *testing.T) {
	suite.Run(t, new(SaleSuite))
}
