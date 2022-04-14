package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"
	salev1 "sale/gen/api/public/v1"
	saleuiv1 "sale/gen/api/ui/v1"
	"sale/internal/app"
	"sale/internal/models"
)

type Server struct {
	saleApp *app.App
}

func NewServer(saleApp *app.App) *Server {
	return &Server{saleApp: saleApp}
}

func (s *Server) Create(ctx context.Context, r *saleuiv1.CreateRequest) (*emptypb.Empty, error) {
	// todo проверка/валидация?
	sale := &models.Sale{
		Name:           r.Name,
		UUID:           uuid.NewString(),
		StartDate:      r.StartDate.AsTime(),
		EndDate:        r.EndDate.AsTime(),
		StopFactor:     uint8(r.StopFactor),
		ContractorType: uint8(r.ContractorType),
		Type:           uint8(r.Type),
		State:          models.SaleIsActive, // todo
		CreatedBy:      "I",
		CreatedAt:      time.Now(),
	}

	err := s.saleApp.CreateSale(sale)
	if err != nil {
		log.Err(err).Msg("error create sale")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) CreateProducts(ctx context.Context, r *saleuiv1.CreateProductsRequest) (*emptypb.Empty, error) {
	/*params, err := excel.ReadByte(importProducts.Params)
	if err != nil {
		log.Err(err).Msg("error read excel file")
	}*/

	params := []*app.Params{}
	for _, p := range r.SaleProducts {
		params = append(params, &app.Params{
			NomenclatureUUID: p.NomenclatureUuid,
			Price:            uint(p.Price),
			MaxCount:         uint(p.MaxCount),
			MaxOrderCount:    uint(p.MaxOrderCount),
			Type:             uint(p.Type),
			IsFeed:           p.IsFeed,
		})
	}
	s.saleApp.AddToFillProductsByRegions(uint(r.SaleId), r.RegionId, params)

	return &emptypb.Empty{}, nil
}

func (s *Server) Sell(ctx context.Context, r *salev1.SellRequest) (*emptypb.Empty, error) {
	// todo транзакция.
	for _, sell := range r.Sell {
		nomenclUUID, err := uuid.Parse(sell.NomenclatureUuid)
		if err != nil {
			return nil, err
		}
		OrderUUID, err := uuid.Parse(sell.OrderUuid)
		if err != nil {
			return nil, err
		}

		// todo проверка что действие допустимо
		// canSelling(selling.Quantity, selling.OrderUuid, selling.NomenclatureUuid)

		err = s.saleApp.Sell(nomenclUUID.String(), uint(sell.Quantity), OrderUUID.String())
		if err != nil {
			return nil, err
		}
	}

	// todo успешный ответ
	return &emptypb.Empty{}, nil
}

func (s *Server) Rest(ctx context.Context, r *salev1.RestRequest) (*salev1.RestResponse, error) {
	rests := []*salev1.RestResponse_Rest{}
	for _, rest := range r.Rest {
		// todo запросы в цикле
		saleProduct, err := s.saleApp.Rest(rest.NomenclatureUuid, uint(rest.RegionId), rest.OrderUuid)
		if err != nil {
			return nil, err
		}

		rests = append(rests, &salev1.RestResponse_Rest{
			NomenclatureUuid: rest.NomenclatureUuid,
			Sale: &salev1.RestResponse_Rest_Sale{
				SaleId: uint32(saleProduct.SaleID),
				//SaleName:   join в запросе
				Quantity: uint32(saleProduct.Available),
				Price:    uint32(saleProduct.Price),
			},
		})
	}

	return &salev1.RestResponse{
		Rest: rests,
	}, nil
}
