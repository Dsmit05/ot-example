package main

import (
	"log"

	"github.com/Dsmit05/ot-example/service-main/config"
	"github.com/Dsmit05/ot-example/service-main/internal/api"
	"github.com/Dsmit05/ot-example/service-main/internal/broker/kfk"
	clir "github.com/Dsmit05/ot-example/service-main/internal/client-service-read"
	"github.com/Dsmit05/ot-example/tracer"
	"github.com/golang/glog"
)

// @title service-main
// @version 1.0.0
// @description starting point server

// @host localhost:8080
func main() {
	log.Println("service-main start")

	cfg := config.Config{
		CollectorURL:   "localhost:4317",
		ServiceReadURL: ":9080",
		ServiceMainURL: ":8080",
	}

	exp, err := tracer.NewExporter(tracer.Config{
		ServiceName:  "service-main",
		CollectorURL: cfg.CollectorURL,
	})
	if err != nil {
		log.Fatalf("tracer.NewExporter error: %v", err)
	}

	producer, err := kfk.NewProducer("localhost:9092", "msgs")
	if err != nil {
		log.Fatalf("kfk.NewProducer() error: %v", err)
	}

	grpcCli, err := clir.NewClient(cfg.ServiceReadURL)
	if err != nil {
		log.Fatalf("clir.NewClient error: %v", err)
	}

	s := api.NewServer(grpcCli, exp.Exp, producer)
	glog.Fatal(s.Start(cfg.ServiceMainURL))
}
