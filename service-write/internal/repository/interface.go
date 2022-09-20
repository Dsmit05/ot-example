package repository

import (
	"context"

	"github.com/Dsmit05/ot-example/service-write/internal/models"
)

type ReposytoryI interface {
	SetMsg(ctx context.Context, msg models.Message) error
}
