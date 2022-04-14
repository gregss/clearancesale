package main

import (
	"context"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	salev1 "sale/gen/api/public/v1"
	"sale/internal/app"
	"sale/internal/app/transport/grpc"
	"sale/internal/config"
	"sale/internal/storage/sql"
)

func main() {
	time.Sleep(5 * time.Second) // todo почему-то при первом старте выпадают ошибки подключения к бд или кролику
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
		log.Err(err).Msg("error new storage")
	}

	handler := salev1.NewSellingServiceServer(server.NewServer(app.NewApp(*cfg, stg)))
	// todo mux.
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Err(err).Msg("error listen and serve rpc server")
	}
}
