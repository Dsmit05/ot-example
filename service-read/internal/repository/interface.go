package repository

import "context"

type ReposytoryI interface {
	GetMsg(ctx context.Context, id int64) (string, error)
}
