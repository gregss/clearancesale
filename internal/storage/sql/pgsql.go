package sql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"sale/internal/app/calculator"
	"sale/internal/models"
)

type Storage struct {
	dsn string
	con *sqlx.DB
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	stg := &Storage{dsn: dsn}

	go func() {
		<-ctx.Done()
		err := stg.Close()
		if err != nil {
			fmt.Println("error") // todo в лог
		}
	}()

	err := stg.Connect()
	if err != nil {
		return nil, err
	}

	err = stg.Ping()
	if err != nil {
		return nil, err
	}

	return stg, nil
}

func (s *Storage) Connect() error {
	var err error
	s.con, err = sqlx.Open("pgx", s.dsn)

	if err != nil {
		log.Fatalf("err %v", err)
	}

	if s.con == nil {
		log.Fatal("nil connection")
	}

	return nil
}

func (s *Storage) Ping() error {
	return s.con.Ping()
}

func (s *Storage) Close() error {
	return s.con.Close()
}

func (s *Storage) CreateSale(sale *models.Sale) error {
	query := `insert into sale(uuid, name, start_date, end_date, stop_factor,
                 contractor_type, type, state, created_by, created_at)
			values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	result, err := s.con.Exec(query,
		sale.UUID,
		sale.Name,
		sale.StartDate,
		sale.EndDate,
		sale.StopFactor,
		sale.ContractorType,
		sale.Type,
		sale.State,
		sale.CreatedBy,
		sale.CreatedAt,
	)
	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("nil result") // странно что может быть nil
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New("0 rows affected")
	}

	// todo подобрать другой драйвер, получать LastInsertId (is not supported by this driver)

	return nil
}

func (s *Storage) EditSale(sale *models.Sale) error {
	_ = sale
	return nil
}

func (s *Storage) AddProduct(p *models.SaleProduct) error {
	query := `insert into sale_product(uuid, sale_id, region_id, nomencl_uuid, price, max_count, max_order_count,
                         stock_available, available, is_feed, status, created_by, created_at, updated_by, updated_at)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	result, err := s.con.Exec(query,
		uuid.New(),
		p.SaleID,
		p.RegionID,
		p.NomenclatureUUID,
		p.Price,
		p.MaxCount,
		p.MaxOrderCount,
		0, // StockAvailable
		0, // Available
		p.IsFeed,
		p.Status,
		uuid.New(), // getUserUUID
		time.Now(),
		uuid.New(), // getUserUUID
		time.Now(),
	)
	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("nil result") // странно что может быть nil
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New("0 rows affected")
	}

	// todo подобрать другой драйвер, получать LastInsertId (is not supported by this drive

	return nil
}

func (s *Storage) EditProduct(p *models.SaleProduct) error {
	_ = p
	return nil
}

func (s *Storage) SaveParams(id uint32, v *calculator.CalcValues) error {
	result, err := s.con.Exec(
		`update sale_product set stock_available = $1, available = $2, status = $3 where id = $4`,
		v.StockAvailable,
		v.Available,
		v.Status,
		id,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count < 1 {
		return err
	}

	return nil
}

func (s *Storage) AddSell(sell *models.Sell) error {
	query := `insert into sell(sale_product_id, order_uuid, quantity, created_at) values($1, $2, $3, $4)`

	result, err := s.con.Exec(query,
		sell.SaleProductID,
		sell.OrderUUID,
		sell.Quantity,
		time.Now(),
	)
	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("nil result") // странно что может быть nil
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return errors.New("0 rows affected")
	}

	return nil
}

// GetPurchasedCount Получаем сумму проданных позиций по всем регионам.
func (s *Storage) GetPurchasedCount(product *models.SaleProduct) (uint, error) {
	var count interface{}
	// todo
	err := s.con.Get(
		&count,
		`select sum(quantity) from sell where sale_product_id = $1`,
		product.ID,
	)

	if count == nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return count.(uint), nil
}

func (s *Storage) GetSaleProduct(nomenclUUID string, regionID uint) (*models.SaleProduct, error) {
	row := s.con.QueryRow(`
		select sale_id, available, price
		from sale_product
		where 
		      nomencl_uuid = $1
		  and region_id = $2
		  and status <> $3	
		`,
		nomenclUUID,
		regionID,
		models.SaleProductIsInactive,
	)

	saleProduct := models.SaleProduct{}
	// todo
	if err := row.Scan(
		&saleProduct.SaleID,
		&saleProduct.Available,
		&saleProduct.Price,
	); err != nil {
		return nil, err
	}

	return &saleProduct, nil
}
