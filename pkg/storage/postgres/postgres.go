package postgres

import (
	"context"
	"time"

	"github.com/pdkonovalov/user-registration-service/pkg/config"
	"github.com/pdkonovalov/user-registration-service/pkg/storage"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	pool *pgxpool.Pool
}

func Init(config *config.Config) (storage.Storage, error) {
	pool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	_, err = pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS email_codes (
		email TEXT PRIMARY KEY,
		code INTEGER NOT NULL,
		exp_time TIMESTAMPTZ NOT NULL
		);
		CREATE INDEX IF NOT EXISTS email_alias ON email_codes USING hash(email);
		CREATE TABLE IF NOT EXISTS users (
		email TEXT PRIMARY KEY,
		name TEXT,
		username TEXT,
		password_hash CHAR(60)
		);
		CREATE INDEX IF NOT EXISTS email_alias ON users USING hash(email);
	`)
	if err != nil {
		return nil, err
	}
	return &postgres{pool}, nil
}

func (db *postgres) Shutdown() error {
	db.pool.Close()
	return nil
}

func (db *postgres) WriteEmailCode(email string, code int, time time.Time) error {
	_, err := db.pool.Exec(context.Background(),
		`INSERT INTO email_codes(email, code, exp_time) VALUES($1, $2, $3);`, email, code, time)
	if err != nil {
		return err
	}
	return nil
}

func (db *postgres) FindEmailCode(email string) (int, time.Time, bool, error) {
	var code int
	var time time.Time
	err := db.pool.QueryRow(context.Background(),
		`SELECT code, exp_time FROM email_codes WHERE email = $1;`, email).Scan(&code, &time)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return code, time, false, nil
		}
		return code, time, false, err
	}
	return code, time, true, nil
}

func (db *postgres) DeleteEmailCode(email string) error {
	_, err := db.pool.Exec(context.Background(),
		`DELETE FROM email_codes WHERE email = $1;`, email)
	if err != nil {
		return err
	}
	return nil
}

func (db *postgres) WriteNewUser(email string, name string, username string, password_hash string) error {
	_, err := db.pool.Exec(context.Background(),
		`INSERT INTO users(email, name, username, password_hash) VALUES($1, $2, $3, $4);`,
		email, name, username, password_hash)
	if err != nil {
		return err
	}
	return nil
}

func (db *postgres) UpdatePassword(email string, password_hash string) error {
	_, err := db.pool.Exec(context.Background(),
		`UPDATE users SET password_hash = $1 WHERE email = $2;`, password_hash, email)
	if err != nil {
		return err
	}
	return nil
}

func (db *postgres) FindUser(email string) (string, string, string, bool, error) {
	var name string
	var username string
	var password_hash string
	err := db.pool.QueryRow(context.Background(),
		`SELECT name, username, password_hash FROM users WHERE email = $1;`,
		email).Scan(&name, &username, &password_hash)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return name, username, password_hash, false, nil
		}
		return name, username, password_hash, false, err
	}
	return name, username, password_hash, true, nil
}
