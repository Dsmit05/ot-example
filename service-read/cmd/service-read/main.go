package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dsmit05/ot-example/service-read/config"
	"github.com/Dsmit05/ot-example/service-read/internal/api"
	"github.com/Dsmit05/ot-example/service-read/internal/repository"
	"github.com/Dsmit05/ot-example/tracer"
)

func main() {
	log.Println("service-read start")

	cfg := config.Config{
		CollectorURL: "localhost:4317",
		GrpcURL:      ":9080",
		HttpURL:      ":9081",
		Database: config.Database{
			Host:     "localhost",
			Port:     5432,
			Table:    "example",
			User:     "postgres",
			Password: "postgres",
		},
	}

	_, err := tracer.NewExporter(tracer.Config{
		ServiceName:  "service-read",
		CollectorURL: cfg.CollectorURL,
	})
	if err != nil {
		log.Fatalf("tracer.NewExporter error: %v", err)
	}

	db, err := repository.NewPostgresRepository(cfg.GetConnStringDB())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	serv := api.NewServer(db)

	go serv.StartGRPC(cfg.GrpcURL)
	go serv.StartHTTP(cfg.HttpURL, cfg.GrpcURL)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	serv.Stop()
}
