package pgsql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"account/library/config"
	"account/library/logger"
	"sync"
	"go.uber.org/zap"
)

type dbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type PgDbInfo struct {
	ServiceName string
	Env         string
	DbConfig    dbConfig
	Conn        *sql.DB
}

var pgInstanceMap = make(map[string]*PgDbInfo, 1)
var lock *sync.Mutex = &sync.Mutex{}

func LoadPgDb(serviceName, env string) *PgDbInfo {

	DbConfig := new(dbConfig)
	err := config.GetConfigMapObj("db", DbConfig, serviceName, env)
	if err == nil {
		PgDbInfo := new(PgDbInfo)
		PgDbInfo.DbConfig = *DbConfig
		PgDbInfo.Env = env
		PgDbInfo.ServiceName = serviceName
		PgDbInfo.InitConnect()
		return PgDbInfo
	}
	return nil
}

func (PgDbInfo *PgDbInfo) InitConnect() {
	dbConf := PgDbInfo.DbConfig
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database)
	if db, e := sql.Open("postgres", connStr); e != nil {
		PgDbInfo.Conn = nil
	} else {
		db.SetMaxIdleConns(20)
		db.SetMaxOpenConns(20)
		PgDbInfo.Conn = db
	}
}

func (PgDbInfo *PgDbInfo) CheckAndReturnConn() *sql.DB {
	if PgDbInfo.Conn == nil {
		lock.Lock()
		defer lock.Unlock()
		if PgDbInfo.Conn == nil {
			PgDbInfo.InitConnect()
		}
	}
	if err := PgDbInfo.Conn.Ping(); err != nil {
		logger.ZapError.Warn("connect pg fail", zap.Error(err))
		return nil
	}
	return PgDbInfo.Conn
}

func (PgDbInfo *PgDbInfo) Clean() {
	if PgDbInfo.Conn != nil {
		logger.ZapTrace.Info("close postage conn")
		errClean := PgDbInfo.Conn.Close()
		if errClean != nil {
			logger.ZapError.Warn("connect pg fail", zap.Error(errClean))
		}
	}
	PgDbInfo.Conn = nil
}
