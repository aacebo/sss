package postgres

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/aacebo/sss/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("name", "postgres")
var conn *sql.DB

func New() *sql.DB {
	if conn == nil {
		url := utils.GetEnv(
			"POSTGRES_CONNECTION_STRING",
			"postgresql://admin:admin@localhost:5432/gladius?sslmode=disable",
		)

		db, err := sql.Open("postgres", url)

		if err != nil {
			panic(err)
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(10)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := db.Ping(); err != nil {
			panic(err)
		}

		conn = db
		driver, err := postgres.WithInstance(db, &postgres.Config{})

		if err != nil {
			panic(err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://postgres/migrations",
			"postgres",
			driver,
		)

		if err != nil {
			panic(err)
		}

		m.Up()
		log.Info("connection established...")
	}

	return conn
}
