package initialize

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go-ecommerce-backend-api.com/global"
)

func InitMysqlC() {
	m := global.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Dbname,
	)

	db, err := sql.Open("mysql", dsn)
	checkErrorPanic(err, "InitMysql initialization error")

	// Set pool cấu hình
	db.SetMaxIdleConns(m.MaxIdleConns)
	db.SetMaxOpenConns(m.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(m.ConnMaxIdleTime) * time.Second)

	// Test kết nối
	err = db.Ping()
	checkErrorPanic(err, "Failed to ping DB")

	global.Logger.Info("✅ MySQL connection pool initialized successfully")
	global.MDBC = db
}
