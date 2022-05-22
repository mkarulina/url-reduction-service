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

func CreateTable(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(
		ctx,
		"CREATE TABLE IF NOT EXISTS urls ("+
			"user_id VARCHAR(255), "+
			"key VARCHAR(255), "+
			"link VARCHAR(255), "+
			"is_deleted BOOLEAN DEFAULT false, "+
			"CONSTRAINT uniq_person_link UNIQUE (user_id, key, link))",
	)
	if err != nil {
		return err
	}

	return nil
}

func AddURLToTable(link *Link) error {
	dbAddress := viper.GetString("DATABASE_DSN")

	db, err := sql.Open("pgx", dbAddress)
	if err != nil {
		return err
	}
	defer db.Close()

	err = CreateTable(db)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insert, err := tx.PrepareContext(
		ctx,
		"INSERT INTO urls (user_id, key, link) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
	)
	if err != nil {
		return err
	}
	defer insert.Close()

	result, err := insert.ExecContext(ctx, link.UserID, link.Key, link.Link)
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

	row := db.QueryRowContext(ctx, "SELECT key, link, is_deleted FROM urls WHERE link = $1 or key = $1", value)
	if err != nil {
		log.Panic(err)
	}

	if row != nil {
		err = row.Scan(&foundLink.Key, &foundLink.Link, &foundLink.IsDeleted)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		}
	}

	return foundLink, nil
}

func GetAllRowsByUserID(userID string) ([]Link, error) {
	var links []Link

	dbAddress := viper.GetString("DATABASE_DSN")

	db, err := sql.Open("pgx", dbAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT key, link FROM urls WHERE user_id = $1 AND is_deleted IS NOT true", userID)
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

func SetIsDeletedFlag(userID string, keys []string) error {
	dbAddress := viper.GetString("DATABASE_DSN")

	db, err := sql.Open("pgx", dbAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insert, err := tx.PrepareContext(
		ctx,
		"UPDATE urls SET is_deleted = true WHERE user_id = $1 AND key = any($2)",
	)
	if err != nil {
		return err
	}
	defer insert.Close()

	_, err = insert.ExecContext(ctx, userID, keys)
	if err != nil {
		return err
	}

	return tx.Commit()
}
