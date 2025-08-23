package initialize

import (
	"database/sql"
	"fmt"
	"time"

	"Taskly.com/m/global"
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
func initPostgresConnection(pg string) *sql.DB {
	// Format connection string

	fmt.Println("dsn pg", pg)
	db, err := sql.Open("postgres", pg)
	checkErrorPanicC(err, fmt.Sprintf("Failed to initialize PostgreSQL "))

	// Optional: Test the connection
	err = db.Ping()
	checkErrorPanicC(err, "Failed to ping PostgreSQL database")

	return db
}

// InitPostgreSQL initializes the PostgreSQL connection and configures pooling
func InitPostgreSQLProd() {
	global.PostgreSQL = initPostgresConnection(global.ENVSetting.Database_url_internal)
	setPostgresPool()
}
func InitPostgreSQLDev() {
	global.PostgreSQL = initPostgresConnection(global.ENVSetting.Database_url_external)
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
