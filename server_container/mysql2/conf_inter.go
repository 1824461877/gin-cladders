package mysql2

import (
	"database/sql"
	"gin-cladder/conf/elite/control"
	"gin-cladder/conf/elite/extend"
	"log"
	"strconv"
	"strings"
	"time"
)

type SMysql struct {
	SMysqlConf `mapstructure:"conf"`
	Conn interface{}
}


type SMysqlConf struct {
	ServiceName string `mapstructure:"service_name"`
	ServiceAddr string `mapstructure:"service_addr"`
	ServicePort int `mapstructure:"service_port"`
	ServiceDatabase string `mapstructure:"service_database"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	MaxConnLifeTime int `mapstructure:"max_conn_life_time"`
	MaxOpenConn int `mapstructure:"max_open_conn"`
	MaxIdleConn int `mapstructure:"max_idle_conn"`
}

func (sec SMysql) Transport(cp *extend.ConnPool) {
	err := control.ParseConfig("server_container/mysql2/conf.toml",&sec)
	if err != nil {
		log.Println("this is err",err)
		return
	}
	sec.ServiceName = strings.ToLower(sec.ServiceName)
	if err := sec.Mysql(); err != nil {
		return
	}
}

func (sec *SMysql) Mysql() error {
	if extend.ConnPoolShareInitServer[sec.ServiceName] != nil {
		return nil
	}
	sec.ServiceDatabase = extend.Parsec(sec.ServiceDatabase)
	connd := sec.User + ":" + sec.Password + "@tcp(" + sec.ServiceAddr + ":" + strconv.Itoa(sec.ServicePort) + ")" + sec.ServiceDatabase
	db, err := sql.Open("mysql",connd)
	if err != nil {
		return nil
	}
	db.SetMaxIdleConns(sec.MaxConnLifeTime)
	db.SetMaxOpenConns(sec.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(sec.MaxConnLifeTime) * time.Second)
	extend.ConnPoolShareInitServer[sec.ServiceName] = db
	return nil
}

func InitMysql2DB(){
	sec := SMysql{}
	extend.ExtendService("Mysql2",sec)
}