package repository

import (
	"context"
	"fmt"

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

const getMsg = `
SELECT msg
FROM msgs
WHERE id = $1;
`

func (o *PostgresRepository) GetMsg(ctx context.Context, id int64) (string, error) {
	var result string

	rows := o.conn.QueryRow(ctx, getMsg, id)

	if err := rows.Scan(&result); err != nil {
		return result, errors.Wrap(err, "rows.Scan() error")
	}

	return result, nil
}

func (o *PostgresRepository) Close() {
	o.conn.Close()
}
