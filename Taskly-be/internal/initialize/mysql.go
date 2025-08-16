package initialize

import (
	"database/sql"
	"fmt"
	"time"

	"Taskly.com/m/global"
	"Taskly.com/m/package/setting"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// checkErrorPanicC handles errors by logging them and triggering a panic
func checkErrorPanicC(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

// initPostgresConnection initializes a PostgreSQL connection using provided config
func initPostgresConnection(pg setting.PostgreSQLSetting) *sql.DB {
	// Format connection string
	dsn := "postgres://%s:%s@%s:%d/%s?sslmode=disable"
	connStr := fmt.Sprintf(dsn, pg.Username, pg.Password, pg.Host, pg.Port, pg.Dbname)

	db, err := sql.Open("postgres", connStr)
	checkErrorPanicC(err, fmt.Sprintf("Failed to initialize PostgreSQL for DB: %s", pg.Dbname))

	// Optional: Test the connection
	err = db.Ping()
	checkErrorPanicC(err, "Failed to ping PostgreSQL database")

	return db
}

// InitPostgreSQL initializes the PostgreSQL connection and configures pooling
func InitPostgreSQL() {
	global.PostgreSQL = initPostgresConnection(global.Config.PostgreSQL)
	setPostgresPool()
}

// setPostgresPool configures the database connection pooling settings
func setPostgresPool() {
	sqlDb := global.PostgreSQL

	// Tuỳ bạn có cần config sâu thêm, ở đây mình để mặc định hoặc bạn có thể thêm MaxOpen/MaxIdle
	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(30 * time.Minute)
}
