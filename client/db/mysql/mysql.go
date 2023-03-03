package mysql

import (
	"fmt"
	"github.com/woodlsy/woodGin/config"
	"github.com/woodlsy/woodGin/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

//type Mysql struct {
//	Db *gorm.DB
//}

func connect() map[string]*gorm.DB {
	if len(config.Configs.Databases) == 0 {
		log.Logger.Error("未配置数据库配置 Databases = ", config.Configs.Databases)
		panic("请先配置数据库配置")
	}

	cons := map[string]*gorm.DB{}
	for _, database := range config.Configs.Databases {
		cons[database.Dbname] = openDb(database)
	}
	return cons
}

func openDb(databases config.Databases) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		databases.UserName,
		databases.Password,
		databases.Host,
		databases.Port,
		databases.Dbname,
		databases.Charset,
	)

	dbLogger := log.NewDbLogger(ormLogger.Info, log.Config{
		SlowThreshold:             time.Second, // 慢 SQL 阈值
		IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  false,       // 禁用彩色打印
	})

	con, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   databases.Prefix, // table name prefix, table for `User` would be `t_users`
			SingularTable: true,             // use singular table name, table for `User` would be `user` with this option enabled
			//NoLowerCase:   true,                              // skip the snake_casing of names
			//NameReplacer: strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to Db name
		},
		Logger: dbLogger,
	})

	if err != nil {
		log.Logger.Error("数据库", databases.Dbname, "连接失败", err, dsn)
		panic("数据库链接失败")
	}
	sqlDB, err := con.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("数据库", databases.Dbname, "连接成功")
	return con
}
