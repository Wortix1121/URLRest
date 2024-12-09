package postgres

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connpool struct {
	dbpool *pgxpool.Pool
}

var (
	pgInst *Connpool
	pgOnce sync.Once
)

// Pool Connections
func NewPG(ctx context.Context, storagePathPg string) (*Connpool, error) {

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, storagePathPg)
		if err != nil {
			return
		}

		pgInst = &Connpool{db}
	})

	return pgInst, nil

}

func (pg *Connpool) Ping(ctx context.Context) error {
	return pg.dbpool.Ping(ctx)
}

func (pg *Connpool) Close() {
	pg.dbpool.Close()
}
