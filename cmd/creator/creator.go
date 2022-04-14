package main

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sale/internal/app"
	"sale/internal/app/creator"
	"sale/internal/config"
	"sale/internal/queue/consumer"
	"sale/internal/storage/sql"
)

// консьюмер, добавляющий позиции в распродажу и делающий расчет остатков.
func main() {
	time.Sleep(5 * time.Second) // todo почему-то при первом старте выпадают ошибки подключения к бд или консьюмеру
	ctx := context.TODO()
	cfg := &config.AppCfg{}
	config.LoadConfig(cfg)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	if cfg.IsDev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	stg, err := sql.New(ctx, cfg.DbURI)
	if err != nil {
		log.Err(err).Msg(cfg.DbURI)
		return
	}

	appSale := app.NewApp(*cfg, stg)

	_, err = consumer.New(cfg.RabbitURI, (&creator.CreateHandler{App: *appSale}).Handle)
	if err != nil {
		log.Err(err).Msg("error consume")
	}
}
