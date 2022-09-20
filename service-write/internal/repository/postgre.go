package repository

import (
	"context"
	"fmt"

	"github.com/Dsmit05/ot-example/service-write/internal/models"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// PostgresRepository repository for accessing Postgres database.
type PostgresRepository struct {
	conn *pgxpool.Pool
}

func NewPostgresRepository(connString string) (*PostgresRepository, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	conn, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return &PostgresRepository{conn: conn}, err
}

const setMsg = `
INSERT INTO msgs (id, msg)
VALUES ($1, $2);
`

func (o *PostgresRepository) SetMsg(ctx context.Context, msg models.Message) error {
	_, err := o.conn.Exec(ctx, setMsg, msg.ID, msg.Msg)
	if err != nil {
		return errors.Wrap(err, "o.conn.Exec() error")
	}

	return nil
}

func (o *PostgresRepository) Close() {
	o.conn.Close()
}
