package db

import (
	"iQuest/config"
	"iQuest/logger"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlDB *gorm.DB
var redisDB *redis.Client

func init() {
	connectMySQL()
	// autoMigrate()
}

// 连接MySQL
func connectMySQL() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}
	// 连接
	var err error

	mysqlDB, err = gorm.Open("mysql", config.Viper.GetString("MYSQL_URL"))

	if err != nil {
		logger.Panic("failed to connect mysql", err)
	}

	mysqlDB.SingularTable(true)
	mysqlDB.LogMode(true)
	mysqlDB.SetLogger(&logger.DBLogger{})
}

// Close 关闭连接
func Close() {
	_ = mysqlDB.Close()
	// if config.Viper.GetBool("REDIS") || config.Viper.GetBool("REDIS_SENTINEL") {
	// 	_ = redisDB.Close()
	// }
}

// Get 获取数据库
func Get() *gorm.DB {
	return mysqlDB
}
