package broker

import (
	"context"

	"github.com/Dsmit05/ot-example/service-main/internal/models"
)

type Producer interface {
	SendMsg(ctx context.Context, msg models.Message) error
}
