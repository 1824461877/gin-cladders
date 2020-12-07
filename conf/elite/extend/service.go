package extend

import (
	"database/sql"
	"fmt"
	"gin-cladder/conf/elite/control"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)


type ServiceExtend struct {
	ServiceExtendConf `mapstructure:"conf"`
	Conn interface{}
}

type ServiceExtendConf struct {
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

func (sec ServiceExtend) Transport(cp *ConnPool) {
	err := cp.Set(sec.ServiceName,PoolFormat{3 * time.Minute,sec.Conn})
	if err != nil{
		fmt.Println("this is transport error ")
		return
	}
}

func init() {
	err := getfilelist("server_container/")
	if err != nil {
		panic("server_container is error")
	}
}

func getfilelist(path string) error {
	if path == "" {
		return errors.New("path is nil")
	}
	fs,_:= ioutil.ReadDir(path)
	for _,fileDir := range fs {
		if fileDir.IsDir(){
			servierpath := path + fileDir.Name() + "/"
			fd,_:= ioutil.ReadDir(servierpath)
			for _,file := range fd {
				if file.Name() == "conf.toml" {
					err := InitServiceExtendConf(servierpath + file.Name())
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}


func InitServiceExtendConf(confpath string) error {
	if confpath == "" {
		return errors.Errorf("conf.toml path is nil ")
	}
	sec := &ServiceExtend{}
	err := control.ParseConfig(confpath,sec)
	if err != nil {
		return nil
	}
	sec.ServiceName = strings.ToLower(sec.ServiceName)
	ServiceFactory(sec)
	return nil
}

// this is service factory
func ServiceFactory(sec *ServiceExtend) {
	if sec.ServiceName == "" {
		return
	}
	switch sec.ServiceName {
	case "mysql":
		sec.Conn = sec.Mysql
		NewPool(sec)
	case "mongodb":
		sec.Conn = sec.Mongodb
		NewPool(sec)
	case "redis":
		sec.Conn = sec.Redis
		NewPool(sec)
	case "etcd":
		sec.Conn = sec.Etcd
		NewPool(sec)
	default:
		sec.Conn = sec.Custom
		NewPool(sec)
	}
}


type Customfunc func(sec interface{})

func (sec *ServiceExtend) Custom(newsec interface{}) error {
	newsec.(PoolConf).Transport(ConnPoolShare)
	return nil
}

// this is mysql factory
// root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true&loc=Asia%2FChongqing
func (sec *ServiceExtend) Mysql() error {
	if ConnPoolShareInitServer[sec.ServiceName] != nil {
		return nil
	}
	sec.ServiceDatabase = Parsec(sec.ServiceDatabase)
	connd := sec.User + ":" + sec.Password + "@tcp(" + sec.ServiceAddr + ":" + strconv.Itoa(sec.ServicePort) + ")" + sec.ServiceDatabase
	db, err := sql.Open("mysql",connd)
	if err != nil {
		return nil
	}
	db.SetMaxIdleConns(sec.MaxConnLifeTime)
	db.SetMaxOpenConns(sec.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(sec.MaxConnLifeTime) * time.Second)
	ConnPoolShareInitServer[sec.ServiceName] = db
	return nil
}

// this is mongodb factory
// mongodb://dbwebauth:dbwebauthfrom@coms.xiaopan233.club:27017/dbweblis
func (sec *ServiceExtend) Mongodb() error {
	if ConnPoolShareInitServer[sec.ServiceName] != nil {
		return nil
	}
	sec.ServiceDatabase = Parsec(sec.ServiceDatabase)
	connd := "mongodb://" + sec.User + ":" + sec.Password + "@" + sec.ServiceAddr + ":" + strconv.Itoa(sec.ServicePort) + sec.ServiceDatabase
	client, err := mongo.NewClient(options.Client().ApplyURI(connd))
	if err != nil {
		return err
	}
	fmt.Println(client)
	return nil
}

func (sec *ServiceExtend) Redis() error {
	fmt.Print("this is reids")
	return nil
}

func (sec *ServiceExtend) Etcd() error {
	fmt.Println("this is etcd")
	return nil
}

func Parsec(dbpath string) string{
	pf := strings.HasPrefix(dbpath,"/")
	sf := strings.HasSuffix(dbpath,"/")
	if !pf {
		dbpath = "/" + dbpath
	}
	if sf {
		dbpath = dbpath[0:len(dbpath)-1]
	}
	return dbpath
}

func ExtendService(sername string,sec interface{}) {
	c := ConnPoolShare
	c.Get(sername).(func(interface{})error)(sec)
}

func Service(sername string) {
	c := ConnPoolShare
	c.Get(sername).(func() error)()
}