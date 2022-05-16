package handlers

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

func (h *handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	dbAddress := viper.GetString("DATABASE_DSN")

	db, err := sql.Open("pgx", dbAddress)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
