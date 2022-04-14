package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"sale/internal/app/calculator"
	"sale/internal/config"
	"sale/internal/models"
	"sale/internal/queue/producer"
	"sale/internal/storage"
)

type CreateMessage struct {
	SaleID     uint
	RegionsIDs []uint32
	Params     Params // todo по ссылке
}

type Params struct {
	NomenclatureUUID string
	Price            uint
	MaxCount         uint
	MaxOrderCount    uint
	Type             uint
	IsFeed           bool
}

type App struct {
	storage storage.Storage
	queue   amqp.Queue
	config  config.AppCfg
}

func NewApp(cfg config.AppCfg, stg storage.Storage) *App {
	return &App{storage: stg, config: cfg}
}

// CreateSale создать распродажу.
func (a *App) CreateSale(sale *models.Sale) error {
	// todo валидация, проверка на пересечение с распродажей того же типа, что-то еще

	// сохранение в бд
	sale.UUID = uuid.NewString()
	sale.CreatedAt = time.Now()
	sale.CreatedBy = uuid.NewString()
	sale.State = models.SaleIsInactive // по дефолту неактивированная
	err := a.storage.CreateSale(sale)
	if err != nil {
		return err
	}

	// выгрузка во внешние системы (kafka)

	return nil
}

// AddToFillProductsByRegions добавляем в очередь.
func (a *App) AddToFillProductsByRegions(saleID uint, regionIDs []uint32, params []*Params) {
	// заполняем товары выбронной распродажи в выбранных регионах
	for _, p := range params {
		message := CreateMessage{
			SaleID:     saleID,
			RegionsIDs: regionIDs,
			Params:     *p,
		}

		b, _ := json.Marshal(&message)

		if err := producer.Publish(a.config.RabbitURI, string(b)); err != nil {
			log.Err(err).Msg("publish to rabbitMQ error")
		}
	}
}

// FillProductsByRegions непосредственно заполняем.
// todo в приложении и продюсеры и консьюмеры, возможно нужно разнести.
func (a *App) FillProductsByRegions(message CreateMessage) {
	// todo, без обновления, пока считаем что загружаем всегда новые данные
	for _, regionID := range message.RegionsIDs {
		if message.Params.Price == 0 {
			continue
		}

		saleProduct := &models.SaleProduct{
			UUID:             uuid.NewString(),
			NomenclatureUUID: message.Params.NomenclatureUUID,
			SaleID:           message.SaleID,
			RegionID:         uint(regionID),
			MaxCount:         message.Params.MaxCount,
			MaxOrderCount:    message.Params.MaxOrderCount,
			Price:            message.Params.Price,
			IsFeed:           message.Params.IsFeed,
			Status:           models.SaleProductIsActive, // дефолтное значение, todo подумать, должно быть неактивно по дефолту
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		v, err := a.calculateValues(saleProduct)
		if err != nil {
			log.Err(err).Msg("error calculate values")
			continue
		}

		saleProduct.StockAvailable = v.StockAvailable
		saleProduct.Available = v.Available
		saleProduct.Status = v.Status

		log.Error().Msg(fmt.Sprintf("%v", v))
		err = a.storage.AddProduct(saleProduct)
		if err != nil {
			log.Err(err).Msg("error add product")
		}
	}
}

// calculateValues расчитываем вычисляемые значения и обновляем в бд.
func (a *App) calculateValues(saleProduct *models.SaleProduct) (*calculator.CalcValues, error) {
	// кидаем ее на проверку лимитов и расчет остатков
	calc, err := calculator.NewCalculator(saleProduct, a.storage)
	if err != nil {
		return nil, err
	}

	calcValues, err := calc.Calculate(saleProduct)
	if err != nil {
		return nil, err
	}

	return calcValues, nil
}

// Sell фиксируем продажу по распродажной цене для данного заказа.
func (a *App) Sell(nomenclUUID string, quantity uint, orderUUID string) error {
	// todo ищем запись saleProduct для данного запроса
	// todo тут будет логика приоритета распродажи, проверка на лимиты через калькулятор
	saleProductID := map[string]uint{
		nomenclUUID: uint(1),
	}[nomenclUUID]

	sell := &models.Sell{
		SaleProductID: saleProductID,
		OrderUUID:     orderUUID,
		Quantity:      quantity,
	}

	if err := a.storage.AddSell(sell); err != nil {
		return err
	}
	return nil
}

// Rest получаем доступные к применению распродажи для данной номенклатуры в данном регионе в данном заказе.
func (a *App) Rest(nomenclUUID string, regionID uint, orderUUID string) (*models.SaleProduct, error) {
	saleProduct, err := a.storage.GetSaleProduct(nomenclUUID, regionID)
	if err != nil {
		return nil, err
	}

	// todo применять ограничение на кол-во распродажных товаров в заказе.
	return saleProduct, nil
}
