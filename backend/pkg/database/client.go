// pkg/database/client.go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"lytemp/internal/domain/models"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // sql.Open("pgx", ...) için

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Client struct {
	db    *gorm.DB
	Debug bool
}

const dsnMySQL string = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s"

// -------------------------
// Public factory functions
// -------------------------

func ConnectMySQL(user, password, host string, port uint16, databaseName, timezone string, debug bool) (*Client, error) {
	dsn := fmt.Sprintf(dsnMySQL, user, password, host, port, databaseName, timezone)
	return connect(mysql.Open(dsn), debug)
}

func ConnectSQLite(path string, debug bool) (*Client, error) {
	return connect(sqlite.Open(path), debug)
}

// ConnectPostgres: hedef DB yoksa otomatik oluşturur ve bağlanır.
func ConnectPostgres(user, password, host string, port uint16, databaseName string, debug bool) (*Client, error) {
	// Önce hedef DB’ye bağlanmayı dene
	cl, err := connectViaPGX(user, password, host, port, databaseName, debug)
	if err == nil {
		return cl, nil
	}
	// "database does not exist" ise DB’yi oluşturup tekrar dene
	if isDatabaseDoesNotExist(err) {
		if err2 := ensureDatabaseExists(user, password, host, port, databaseName); err2 != nil {
			return nil, fmt.Errorf("create database failed: %w", err2)
		}
		return connectViaPGX(user, password, host, port, databaseName, debug)
	}
	return nil, fmt.Errorf("could not connect to database: %w", err)
}

// -------------------------
// Internals (Postgres)
// -------------------------

func connectViaPGX(user, password, host string, port uint16, databaseName string, debug bool) (*Client, error) {
	// Lokal geliştirme için sslmode=disable; application_name loglarda iş görür
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable application_name=lytemp",
		host, port, user, databaseName, password,
	)

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// Bağlantı hatalarını (ör. 3D000) erken yakalamak için pingle
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return connect(postgres.New(postgres.Config{Conn: sqlDB}), debug)
}

func isDatabaseDoesNotExist(err error) bool {
	if err == nil {
		return false
	}
	// GORM/pgx hata sarmalamasında SQLSTATE metne yansıyor: 3D000
	msg := err.Error()
	msgLower := strings.ToLower(msg)
	return strings.Contains(msg, "SQLSTATE 3D000") ||
		(strings.Contains(msgLower, "database") && strings.Contains(msgLower, "does not exist"))
}

// postgres db’sine bağlanıp hedef DB yoksa CREATE DATABASE ile üret
func ensureDatabaseExists(user, password, host string, port uint16, databaseName string) error {
	adminDSN := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=postgres password=%s sslmode=disable application_name=lytemp",
		host, port, user, password,
	)
	db, err := sql.Open("pgx", adminDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	var exists bool
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);`, databaseName).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}

	// OWNER’ı bağlanan kullanıcı yapıyoruz (ihtiyaca göre değiştir)
	create := fmt.Sprintf(`CREATE DATABASE "%s" WITH OWNER "%s" TEMPLATE template1 ENCODING 'UTF8'`, databaseName, user)
	if _, err := db.Exec(create); err != nil {
		// duplicate_database (42P04) olursa görmezden gel (yarış koşulu vb.)
		if !isDuplicateDatabaseErr(err) {
			return err
		}
	}
	return nil
}

func isDuplicateDatabaseErr(err error) bool {
	if err == nil {
		return false
	}
	e := err.Error()
	return strings.Contains(e, "SQLSTATE 42P04") || // duplicate_database
		strings.Contains(strings.ToLower(e), "already exists")
}

// -------------------------
// Shared helpers
// -------------------------

func connect(dialector gorm.Dialector, debug bool) (*Client, error) {
	var c Client

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %s", err.Error())
	}

	c.db = db
	c.Debug = debug

	if c.Debug {
		c.db = c.db.Debug()
	}

	return &c, nil
}

func (c *Client) HealthCheck(ctx context.Context) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database is not reachable: %s", err.Error())
	}

	return nil
}

func (c *Client) Migrate() error {
	/*
		SeedEnums(c.db)
	*/
	err := c.db.AutoMigrate(
		&models.Content{},
		&models.ProviderSync{},
	)
	if err != nil {
		return fmt.Errorf("could not migrate: %s", err.Error())
	}
	return nil
}

func (c *Client) Get() *gorm.DB {
	return c.db
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *Client) Begin() *Client {
	return &Client{
		db:    c.db.Begin(),
		Debug: c.Debug,
	}
}

func (c *Client) Commit() error {
	return c.db.Commit().Error
}

func (c *Client) Rollback() error {
	return c.db.Rollback().Error
}

func (c *Client) WithTX(tx *Client) *Client {
	return &Client{
		db:    tx.db,
		Debug: c.Debug,
	}
}
