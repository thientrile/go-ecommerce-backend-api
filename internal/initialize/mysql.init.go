package initialize

import (
	"fmt"
	"time"

	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/internal/po"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func checkErrorPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func InitMysql() {
	m := global.Config.Mysql

	dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.Username, m.Password, m.Host, m.Port, m.Dbname)
	db, err := gorm.Open(mysql.Open(s), &gorm.Config{
		SkipDefaultTransaction: false,
	})
	checkErrorPanic(err, "InitMysql  initialization error")
	global.Logger.Info("InitMysql initialization successful")
	global.MDB = db

	//set Pool
	SetPool()
	//migrate tables
	migrateTables()
	global.Logger.Info("Mysql connection pool and tables migration completed successfully")
}

func SetPool() {
	m := global.Config.Mysql
	sqlDb, err := global.MDB.DB()
	if err != nil {
		fmt.Printf("Mysql error:: %s", err)
	}
	sqlDb.SetConnMaxIdleTime(time.Duration(m.MaxIdleConns))
	sqlDb.SetMaxOpenConns(m.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))

}

func migrateTables() {
	err := global.MDB.AutoMigrate(
		&po.User{},
		&po.Role{},
	)
	if err != nil {
		fmt.Printf("Migrate Tables err:: %s", err)
	}

}
