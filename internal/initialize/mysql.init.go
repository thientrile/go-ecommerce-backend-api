package initialize

import (
	"fmt"
	"time"

	"go-ecommerce-backend-api.com/global"
	// "go-ecommerce-backend-api.com/internal/model"
	"go-ecommerce-backend-api.com/internal/po"

	// "go-ecommerce-backend-api.com/internal/po"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func checkErrorPanic(err error, errString string) {
	if err != nil {

		global.Logger.Fatal(errString, zap.Error(err))
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
	// generate table dao
	// GenTableDAO()
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

func GenTableDAO() {
	// gen get apply for file model
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(global.MDB) // reuse your gorm db
	g.GenerateModel("go_crm_user")
	// Generate basic type-safe DAO API for struct `model.User` following conventions
	// g.ApplyBasic(model.User{})

	// // Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

	// Generate the code
	g.Execute()
}

func migrateTables() {
	err := global.MDB.AutoMigrate(
		&po.User{},
		&po.Role{},
		// &model.GoCrmUserV2{},
	)
	if err != nil {
		fmt.Printf("Migrate Tables err:: %s", err)
	}

}
