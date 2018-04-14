package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
)

func GetDatabase(spec config.Specification) *gorm.DB {
	logger := config.GetLogger()

	var db *gorm.DB
	var err error

	// Use sqlite3 for `TEST` env
	if config.IsTestEnv(spec) {
		uuidString := uuid.NewV4().String()
		filename := fmt.Sprintf("/tmp/go-ref_gateway_%s.db", uuidString)
		logger.Infow("Use sqlite3 database", "env", spec.Env)
		db, err = gorm.Open("sqlite3", filename)
	} else {
		logger.Infow("Use mysql database", "env", spec.Env)
		dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			spec.MysqlUserName, spec.MysqlPassword, spec.MysqlHost, spec.MysqlPort, spec.MysqlDatabase)
		db, err = gorm.Open("mysql", dbConnString)
	}

	if err != nil {
		logger.Fatalw("Failed to connect DB", "error", err, "env", spec.Env)
	}

	// migration
	logger.Info("Migrating tables")
	db.SingularTable(true)

	option := "ENGINE=InnoDB"
	if config.IsTestEnv(spec) {
		option = ""
	}

	db.Set("gorm:table_options", option).AutoMigrate(&Access{})
	db.Set("gorm:table_options", option).AutoMigrate(&Session{})

	return db
}
