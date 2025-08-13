package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

// 日志文件夹 /var/logs/card
var logDir = "/var/log/card"

var Logger zerolog.Logger

var logConfig = &lumberjack.Logger{
	Filename:   fmt.Sprintf("%s/log_all.log", logDir), // 日志文件的位置
	MaxSize:    10,                                    // 文件最大尺寸（以MB为单位）
	MaxBackups: 20,                                    // 保留的最大旧文件数量
	MaxAge:     14,                                    // 保留旧文件的最大天数
	Compress:   true,                                  // 是否压缩/归档旧文件
	LocalTime:  true,                                  // 使用本地时间创建时间戳
}

func init() {
	zerolog.TimeFieldFormat = time.DateTime

	// 创建log目录
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("Mkdir failed, err:", err)
		return
	}
	zerolog.CallerSkipFrameCount = 3 // 设置调用函数的层数，默认为2，这里设置为3，即调用函数为main.main()时，层数为3
	Logger = zerolog.New(logConfig).With().Caller().Timestamp().Logger()
	//设置日志等级
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// 设置日志格式
	Logger = Logger.Output(zerolog.ConsoleWriter{
		Out:        io.MultiWriter(os.Stdout, logConfig),
		NoColor:    true,
		TimeFormat: time.DateTime, // 设置时间格式
	})
}

func Info(msg string, args ...any) {
	Logger.Info().Msgf(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error().Msgf(msg, args...)
}
func Err(err error, msg string, args ...any) {
	Logger.Err(err).Msgf(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug().Msgf(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn().Msgf(msg, args...)
}

func Fatal(msg string, args ...any) {
	Logger.Fatal().Msgf(msg, args...)
}
