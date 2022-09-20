package main

import (
	"log"

	"github.com/Dsmit05/ot-example/service-write/internal/broker/kfk"
	"github.com/Dsmit05/ot-example/service-write/internal/config"
	"github.com/Dsmit05/ot-example/service-write/internal/repository"
	"github.com/Dsmit05/ot-example/tracer"
)

func main() {
	log.Println("service-write with kafka start")

	cfg := config.Config{
		CollectorURL: "localhost:4317",
		Database: config.Database{
			Host:     "localhost",
			Port:     5432,
			Table:    "example",
			User:     "postgres",
			Password: "postgres",
		},
	}

	_, err := tracer.NewExporter(tracer.Config{
		ServiceName:  "service-write",
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

	// ToDo: it is necessary to separate the entities of the broker and the database
	br, err := kfk.NewConsumer("localhost:9092", "msgs", db)
	if err != nil {
		log.Fatal(err)
	}

	br.Process()
}
