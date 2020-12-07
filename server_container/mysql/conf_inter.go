package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gin-cladder/conf/elite/extend"
)


type MysqlDB struct {
	SqlDB *sql.DB
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func InitMysqlDB(){
	extend.Service("Mysql")
	server := extend.ConnPoolShareInitServer
	sqldb := server["mysql"].(*sql.DB)
	if err := sqldb.Ping(); err != nil {
		return
	}
	s := MysqlDB{sqldb}
	var userslice []User

	s.Insert("comscs",[]string{"username","password"},"coms","7725032")
	s.Select("comscs",[]string{"username","password"},nil,[]map[string]interface{}{
		{"username":"asc"},
	},func(query *sql.Rows) error {
		for query.Next() {
			var user User
			query.Scan(&user.Username,&user.Password)
			userslice = append(userslice,user)
		}
		return nil
	})
	b , err := json.Marshal(userslice)
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("data:%v",string(b)))
}