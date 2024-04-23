package log

import (
	"context"
	"errors"
	"fmt"
	"github.com/woodlsy/woodGin/config"
	ormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type Int ormLogger.Interface

type Config struct {
	SlowSqlThreshold          int64 // 慢 SQL 阈值，单位：毫秒
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	LogLevel                  ormLogger.LogLevel
}

type DbLogger struct {
	LogLevel ormLogger.LogLevel
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewDbLogger(level ormLogger.LogLevel, config Config) *DbLogger {

	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = ormLogger.Green + "%s\n" + ormLogger.Reset + ormLogger.Green + "[info] " + ormLogger.Reset
		warnStr = ormLogger.BlueBold + "%s\n" + ormLogger.Reset + ormLogger.Magenta + "[warn] " + ormLogger.Reset
		errStr = ormLogger.Magenta + "%s\n" + ormLogger.Reset + ormLogger.Red + "[error] " + ormLogger.Reset
		traceStr = ormLogger.Green + "%s\n" + ormLogger.Reset + ormLogger.Yellow + "[%.3fms] " + ormLogger.BlueBold + "[rows:%v]" + ormLogger.Reset + " %s"
		traceWarnStr = ormLogger.Green + "%s " + ormLogger.Yellow + "%s\n" + ormLogger.Reset + ormLogger.RedBold + "[%.3fms] " + ormLogger.Yellow + "[rows:%v]" + ormLogger.Magenta + " %s" + ormLogger.Reset
		traceErrStr = ormLogger.RedBold + "%s " + ormLogger.MagentaBold + "%s\n" + ormLogger.Reset + ormLogger.Yellow + "[%.3fms] " + ormLogger.BlueBold + "[rows:%v]" + ormLogger.Reset + " %s"
	}

	return &DbLogger{
		LogLevel:     level,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *DbLogger) LogMode(level ormLogger.LogLevel) ormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l DbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ormLogger.Info {
		Logger.Info(l.infoStr+msg, utils.FileWithLineNum(), data)
		//l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l DbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ormLogger.Warn {
		Logger.Warn(l.warnStr+msg, utils.FileWithLineNum(), data)
	}
}

// Error print error messages
func (l DbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ormLogger.Error {
		Logger.Error(l.errStr+msg, utils.FileWithLineNum(), data)
	}
}

// Trace print sql message
func (l DbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	//if l.LogLevel <= ormLogger.Silent {
	//	return
	//}
	sql, rows := fc()
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= ormLogger.Error && (!errors.Is(err, ormLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):

		if rows == -1 {
			Logger.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			Logger.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.SlowSqlThreshold != 0 && elapsed.Milliseconds() > l.SlowSqlThreshold && l.LogLevel >= ormLogger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL > %vms", l.SlowSqlThreshold)
		if rows == -1 {
			Logger.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			Logger.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case config.Configs.App.PSql:
		if rows == -1 {
			Logger.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			Logger.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
