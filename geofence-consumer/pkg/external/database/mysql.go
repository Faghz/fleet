package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type MySQLConnection struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// CreateMySQLConnection returns a MySQL database connection instance
func CreateMySQLConnection(connConfig MySQLConnection, logger *zap.Logger) *sqlx.DB {
	// MySQL DSN format: user:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		connConfig.User, connConfig.Password, connConfig.Host, connConfig.Port, connConfig.DbName)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to MySQL database", zap.Error(err))
		return nil
	}

	// Configure connection pool
	db.SetMaxOpenConns(connConfig.MaxOpenConns)
	db.SetMaxIdleConns(connConfig.MaxIdleConns)
	db.SetConnMaxLifetime(connConfig.ConnMaxLifetime)

	// Test the connection
	err = db.Ping()
	if err != nil {
		logger.Fatal("Failed to ping MySQL database", zap.Error(err))
		return nil
	}

	logger.Info("MySQL Connected successfully")

	return db
}
