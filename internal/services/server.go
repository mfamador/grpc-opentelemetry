// Package services contains the business logic for Posting messages
package services //nolint
import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Server exposes the interface to do operations on the Packets entity
type Server interface {
	Post(ctx context.Context, message string) error
}

type server struct{}

// NewService creates a Packets service
func NewService() Server {
	return &server{}
}

func openDB(dsn string) (*sql.DB, error) {
	driverName, err := otelsql.Register("postgres",
		otelsql.AllowRoot(),
		otelsql.TraceAll(),
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		return nil, err
	}

	// Connect to a Postgres database using the postgres driver wrapper.
	return sql.Open(driverName, dsn)
}

// Post posts messages
func (s *server) Post(ctx context.Context, message string) error {
	log.Info().Msg(message)

	db, err := openDB(fmt.Sprintf("host=localhost port=5432 user=postgres password=postgres sslmode=disable"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.QueryContext(ctx, "select 1")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	for res.Next() {
		var id int
		if err := res.Scan(&id); err != nil {
			panic(err)
		}
		fmt.Println(id)
	}
	return nil
}
