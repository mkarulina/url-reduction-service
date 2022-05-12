package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"log"
	"time"
)

func AddURLToTable(link *Link) error {
	dbAddress := viper.GetString("DATABASE_DSN")

	db, err := sql.Open("pgx", dbAddress)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(
		ctx,
		"CREATE TABLE IF NOT EXISTS urls (key varchar(255) UNIQUE, link varchar(255) UNIQUE)",
	)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insert, err := tx.PrepareContext(
		ctx,
		"INSERT INTO urls (key, link) VALUES ($1, $2) ON CONFLICT (link) DO NOTHING",
	)
	if err != nil {
		return err
	}
	defer insert.Close()

	result, err := insert.ExecContext(ctx, link.Key, link.Link)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected < 1 {
		return errors.New(pgerrcode.UniqueViolation)
	}

	return tx.Commit()
}

func FindValueInDB(value string) (Link, error) {
	var foundLink Link

	dbAddress := viper.GetString("DATABASE_DSN")

	db, err := sql.Open("pgx", dbAddress)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, "SELECT * FROM urls WHERE link = $1 or key = $1", value)
	if err != nil {
		log.Panic(err)
	}

	if row != nil {
		err = row.Scan(&foundLink.Key, &foundLink.Link)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		}
	}

	return foundLink, nil
}

func GetAllRows() ([]Link, error) {
	var links []Link

	dbAddress := viper.GetString("DATABASE_DSN")

	db, err := sql.Open("pgx", dbAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT * FROM urls")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		var l Link
		err = rows.Scan(&l.Key, &l.Link)
		if err != nil {
			return nil, err
		}

		links = append(links, l)
	}

	return links, nil
}
