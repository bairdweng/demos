package db

import (
	"iQuest/config"
	"iQuest/logger"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlDB *gorm.DB
var redisDB *redis.Client

func init() {
	connectMySQL()
	// autoMigrate()
	connectRedis()
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

// 连接Redis
func connectRedis() {

	if config.Viper.GetBool("REDIS") {
		redisDB = redis.NewClient(&redis.Options{
			Addr:         config.Viper.GetString("REDIS_HOST") + ":" + config.Viper.GetString("REDIS_PORT"),
			Password:     config.Viper.GetString("REDIS_PWD"),
			DB:           config.Viper.GetInt("REDIS_DB"),
			MaxRetries:   5,
			DialTimeout:  time.Second * 15,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		})
		//defer redisDB.Close()
		_, err := redisDB.Ping().Result()
		if err != nil {
			logger.Panic("failed to connect redis", err)
		}
	} else if config.Viper.GetBool("REDIS_SENTINEL") {
		addrs := []string{
			config.Viper.GetString("REDIS_SENTINEL_0_HOST") + ":" + config.Viper.GetString("REDIS_SENTINEL_0_PORT"),
			config.Viper.GetString("REDIS_SENTINEL_1_HOST") + ":" + config.Viper.GetString("REDIS_SENTINEL_1_PORT"),
			config.Viper.GetString("REDIS_SENTINEL_2_HOST") + ":" + config.Viper.GetString("REDIS_SENTINEL_2_PORT"),
		}
		redisDB = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    config.Viper.GetString("REDIS_SENTINEL_SERVICE"),
			SentinelAddrs: addrs,
			Password:      config.Viper.GetString("REDIS_SENTINEL_PASSWORD"),
			DB:            config.Viper.GetInt("REDIS_SENTINEL_DB"),
		})
		//defer redisDB.Close()
		_, err := redisDB.Ping().Result()
		if err != nil {
			logger.Panic("failed to connect redis", err)
		}
	}

}

func autoMigrate() {
	/*mysqlDB.AutoMigrate(
		&task.Task{},
		&task.Member{},
		&task.CompanyUserBlacklist{},
		&task.CompanyUserKpi{},
		&task.KpiUploadLog{},
		&commonly.CommonlyUsed{},
	)*/
}

// Close 关闭连接
func Close() {
	_ = mysqlDB.Close()
	if config.Viper.GetBool("REDIS") || config.Viper.GetBool("REDIS_SENTINEL") {
		_ = redisDB.Close()
	}
}

// Get 获取数据库
func Get() *gorm.DB {
	return mysqlDB
}

// Redis 返回redis客户端
func Redis() *redis.Client {
	return redisDB
}
