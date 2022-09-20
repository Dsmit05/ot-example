.PHONY: main-with-kafka
main-with-kafka:
	go run github.com/Dsmit05/ot-example/service-main/cmd/with-kafka

.PHONY: main-with-rmq
main-with-rmq:
	go run github.com/Dsmit05/ot-example/service-main/cmd/with-rmq

.PHONY: write-with-kafka
write-with-kafka:
	go run github.com/Dsmit05/ot-example/service-write/cmd/with-kafka

.PHONY: write-with-rmq
write-with-rmq:
	go run github.com/Dsmit05/ot-example/service-write/cmd/with-rmq

.PHONY: read
read:
	go run github.com/Dsmit05/ot-example/service-read/cmd/service-read