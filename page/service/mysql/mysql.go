package mysql

import(
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
    "sync"
	"go.uber.org/zap"
	"account/library/config"
	"account/library/logger"
	"fmt"
)

type dbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type MysqlDbInfo struct {
	Env         string
	DbConfig    dbConfig
	Conn        *gorm.DB
}

var mysqlInstanceMap = make(map[string]*MysqlDbInfo, 1)
var lock *sync.Mutex = &sync.Mutex{}

func LoadMysqlConn(env string) *MysqlDbInfo {

	DbConfig := new(dbConfig)
	err := config.GetConfigMapObj("db", DbConfig, "mysql", env)
	if err == nil {
		MysqlDbInfo := new(MysqlDbInfo)
		MysqlDbInfo.DbConfig = *DbConfig
		MysqlDbInfo.Env = env
		MysqlDbInfo.InitConnect()
		return MysqlDbInfo
	}
	return nil
}

func (MysqlDbInfo *MysqlDbInfo) InitConnect() {
	dbConf := MysqlDbInfo.DbConfig
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database)
	if db, e := gorm.Open("mysql", connStr); e != nil {
		logger.ZapError.Warn("connect mysql service fail", zap.Error(e))
		MysqlDbInfo.Conn = nil
	} else {
		logger.ZapTrace.Info("connect mysql service success", zap.String("conn", connStr))
		db.DB().SetMaxIdleConns(20)
		db.DB().SetMaxOpenConns(20)
		MysqlDbInfo.Conn = db
	}
}

func (MysqlDbInfo *MysqlDbInfo) CheckAndReturnConn() *gorm.DB {
	if MysqlDbInfo.Conn == nil {
		lock.Lock()
		defer lock.Unlock()
		if MysqlDbInfo.Conn == nil {
			MysqlDbInfo.InitConnect()
		}
	}
	if err := MysqlDbInfo.Conn.DB().Ping(); err != nil {
		logger.ZapError.Warn("connect mysql fail", zap.Error(err))
		MysqlDbInfo.Clean()
		return nil
	}
	return MysqlDbInfo.Conn
}

func (MysqlDbInfo *MysqlDbInfo) Clean() {
	if MysqlDbInfo.Conn != nil {
		logger.ZapTrace.Info("close mysql conn")
		errClean := MysqlDbInfo.Conn.Close()
		if errClean != nil {
			logger.ZapError.Warn("connect mysql fail", zap.Error(errClean))
		}
	}
	MysqlDbInfo.Conn = nil
}
