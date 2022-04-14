package main

import (
	"context"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	saleuiv1 "sale/gen/api/ui/v1"
	"sale/internal/app"
	server "sale/internal/app/transport/grpc"
	"sale/internal/config"
	"sale/internal/storage/sql"
)

// const tsLayout = "2006-01-02T15:04:05"

// приложение, ui сервер.
func main() {
	time.Sleep(5 * time.Second)
	ctx := context.TODO()

	cfg := &config.AppCfg{}
	config.LoadConfig(cfg)

	// todo вынести в утилиты
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	if cfg.IsDev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	stg, err := sql.New(ctx, cfg.DbURI)
	if err != nil {
		log.Err(err).Msg("error new storage")
	}

	handler := saleuiv1.NewSaleServiceServer(server.NewServer(app.NewApp(*cfg, stg)))
	// todo mux.
	err = http.ListenAndServe(":8081", handler)
	if err != nil {
		log.Err(err).Msg("error listen and serve rpc server")
	}
}
