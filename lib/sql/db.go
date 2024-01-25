package sql

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/betam/glb/lib/try"
)

type Db interface {
	Connect() *sqlx.DB
	Timeout() time.Duration
}

func NewDb(
	driver string,
	dsn string,
	timeout time.Duration,
) *db {
	return &db{
		driver:  driver,
		dsn:     dsn,
		timeout: timeout,
	}
}

type db struct {
	driver  string
	dsn     string
	timeout time.Duration
	db      *sqlx.DB
}

func (c *db) Connect() *sqlx.DB {
	if c.db == nil || c.db.Ping() != nil {
		c.db = try.Throw(sqlx.Connect(c.driver, c.dsn))
	}
	return c.db
}

func (c *db) Timeout() time.Duration {
	return c.timeout
}

func (c *db) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
