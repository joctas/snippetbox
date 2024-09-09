package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/joctas/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql" // _ because it will never be directly used
)

// dependency injection
type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	// Flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	dns := flag.String("dns", "web:123456@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dns)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	// server
	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
