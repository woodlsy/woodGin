package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/woodlsy/woodGin/config"
	"github.com/woodlsy/woodGin/helper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var Logger *zap.SugaredLogger

var logFilePath string

func Enabled() {

	var stdoutLevel zapcore.Level
	if config.Configs.App.Debug {
		stdoutLevel = zapcore.DebugLevel
	} else {
		stdoutLevel = zapcore.ErrorLevel
	}

	getFilePathByConfig()

	infoFile := openLogFile("info")
	errorFile := openLogFile("error")
	warnFile := openLogFile("warn")
	debugFile := openLogFile("debug")

	encoder := getEncoder()
	consoleEncoder := getConsoleEncoder()
	newCore := zapcore.NewTee(
		// 写入文件
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),
			infoFile,
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level == zap.InfoLevel
			}),
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),
			errorFile,
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level == zap.ErrorLevel
			}),
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),
			warnFile,
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level == zap.WarnLevel
			}),
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),
			debugFile,
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level == zap.DebugLevel
			}),
		),
		// 写入控制台
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), stdoutLevel),
	)

	//core := zapcore.NewCore(
	//	zapcore.NewConsoleEncoder(encoder),
	//	file,
	//	zapcore.InfoLevel,
	//)

	l := zap.New(newCore, zap.AddCaller())
	zap.ReplaceGlobals(l)
	Logger = l.Sugar()
	fmt.Println(helper.Now(), "日志模块加载成功")
}

func getFilePathByConfig() {
	dir, _ := filepath.Split(config.Configs.Log.FilePath)
	if string(dir[len(dir)-1]) != "/" {
		dir = helper.Join("", dir, "/")
	}
	logFilePath = dir
}

func openLogFile(level string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   helper.Join("", logFilePath, level, ".log"),
		MaxSize:    config.Configs.Log.MaxSize,
		MaxBackups: config.Configs.Log.MaxBackups,
		MaxAge:     config.Configs.Log.MaxAge,
		Compress:   config.Configs.Log.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GetConsoleEncoder 输出日志到控制台
func getConsoleEncoder() zapcore.Encoder {

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getEncoder
// @Description: 日志格式编码
// @return zapcore.EncoderConfig
func getEncoder() zapcore.EncoderConfig {
	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "Llogger",
		CallerKey:      "caller_line",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeDuration: zapcore.SecondsDurationEncoder,

		//MessageKey:     "message",
		//TimeKey:        "time",
		//CallerKey:      "caller_line",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   customTimeEncoder,
		EncodeCaller: cEncodeCaller,
	}
}

// cEncodeCaller 自定义行号显示
func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

// GinLogger
// @Description: gin的日志注入
// @return gin.HandlerFunc
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Debug(fmt.Sprintf("[url] [%s] %s", cost, path),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			//zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			//zap.Duration("cost", cost),
		)
	}
}

func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func StackTrace(msg interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	Logger.Error(msg)
	Logger.Error(string(buf[:n]))
}
