package logger

import (
	"iQuest/config"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {

	fileOut := &lumberjack.Logger{
		Filename:   config.Viper.GetString("LOG_FILENAME"),
		MaxSize:    config.Viper.GetInt("LOG_MAX_SIZE"),    // megabytes
		MaxBackups: config.Viper.GetInt("LOG_MAX_BACKUPS"), // MaxBackups
		MaxAge:     config.Viper.GetInt("LOG_MAX_AGE"),     // days
		LocalTime:  true,                                   // 这个需要设置, 不然日志文件的名字就是UTC时间
	}

	var out io.Writer
	if config.Viper.GetBool("LOG_CONSOLE") {
		out = io.MultiWriter(os.Stdout, fileOut)
	} else {
		out = fileOut
	}

	if config.Viper.GetBool("LOG_JSON") {
		log.Logger = log.Output(out)
	} else {
		writer := zerolog.NewConsoleWriter()
		writer.NoColor = true
		writer.TimeFormat = time.RFC3339
		writer.FormatLevel = func(i interface{}) string {
			if ll, ok := i.(string); ok {
				return strings.ToUpper("[" + ll + "]")
			}
			return "[???]" // level为空
		}
		writer.Out = out

		log.Logger = log.Output(writer)

	}
	log.Level(zerolog.Level(config.Viper.GetUint("LOG_LEVEL")))
}

// Panic 输出错误，中断程序
func Panic(message string, err error) {
	log.Panic().Str("err", err.Error()).Msg(message)
}
