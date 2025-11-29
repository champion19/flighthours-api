package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/champion19/flighthours-api/config"
	loggerPkg "github.com/champion19/flighthours-api/platform/logger"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB(dbConfig config.Database, logger loggerPkg.Logger) (*sql.DB, error) {
	var dsn string

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	if dbConfig.SSL != "" {
		dsn += "&tls=" + dbConfig.SSL
	}

	db, err := sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		logger.Error(loggerPkg.LogDBConnectionError,
			"error", err,
			"host", dbConfig.Host,
			"database", dbConfig.Name)
		return nil, fmt.Errorf("error to connect to database: %w", err)
	}

	logger.Debug(loggerPkg.LogDBConnectionEstablished,
		"max_open_conns", dbConfig.MaxOpenConns,
		"max_idle_conns", dbConfig.MaxIdleConns,
		"conn_max_lifetime", dbConfig.ConnMaxLifetime,
		"conn_max_idle_time", dbConfig.ConnMaxIdleTime,
	)

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime))
	db.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime))

	err = db.Ping()
	if err != nil {
		logger.Error(loggerPkg.LogAppDatabasePingError,
			"error", err,
			"host", dbConfig.Host,
			"database", dbConfig.Name)
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	logger.Success(loggerPkg.LogAppDatabasePingOK,
		"host", dbConfig.Host,
		"database", dbConfig.Name)
	return db, nil
}
