package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// parse slice function string
func parseSliceString(args []string) (string,string) {
	strtag := strings.Join(args,",")
	var strvalueslice []string
	for range args {
		strvalueslice = append(strvalueslice,"?")
	}
	strvalue := strings.Join(strvalueslice,",")
	return strtag,strvalue
}

// parse slice function screen
func parseSliceTagAndValue(args []string) string {
	if args == nil {
		return ""
	}
	var strvalueslice []string
	for _,v := range args {
		strvalueslice = append(strvalueslice,fmt.Sprintf("%v = ?",v))
	}
	strvalue := strings.Join(strvalueslice,",")
	return strvalue
}

// parse slice function screen
func parseSliceTagAndValueSetScreen(screen string,args []string) string {
	if args == nil {
		return ""
	}
	var strvalueslice []string
	for _,v := range args {
		strvalueslice = append(strvalueslice,fmt.Sprintf("%v = ?",v))
	}
	strvalue := strings.Join(strvalueslice,",")
	strvalue = screen + " " + strvalue
	return strvalue
}

// parse map function screen
func parseMapTagAndValueSetScreen(screen string,args []map[string]interface{}) string {
	var strval = screen
	for i,v := range args{
		b , err := json.Marshal(v)
		if err != nil {
			return ""
		}
		str := string(b)
		str = strings.Replace(str,"{","",-1)
		str = strings.Replace(str,"}","",-1)
		str = strings.Replace(str,":"," ",-1)
		if (len(args) - 1)  == i {
			strval += " " + str
		} else {
			strval += " " + str + " ,"
		}
	}
	return strval
}

// this is mysql insert
func (db *MysqlDB) Insert(tablename string,field []string,values ...interface{}) (sql.Result,error) {
	DBSet := db.SqlDB
	strtag,strvalue := parseSliceString(field)
	prepare := fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v)",tablename,strtag,strvalue)
	stmtIns , err := DBSet.Prepare(prepare)
	defer stmtIns.Close()
	if err != nil {
		return nil,fmt.Errorf("db prepare is error : %v",err)
	}
	result, err := stmtIns.Exec(values...)
	if err != nil {
		return nil,fmt.Errorf("db exec is error : %v",err)
	}
	stmtIns.Close()
	return result,nil
}

// this is mysql select
func (db *MysqlDB) Select(tablename string,field []string,screen []string,order []map[string]interface{},funcs func(*sql.Rows) error) ([]string,error) {
	DBSet := db.SqlDB
	strtag,_ := parseSliceString(field)
	screen_value := parseSliceTagAndValueSetScreen("where",screen)
	order_value := parseMapTagAndValueSetScreen("order by",order)

	prepare := fmt.Sprintf("SELECT %v FROM %v %v %v",strtag,tablename,screen_value,order_value)

	query , err := DBSet.Query(prepare)
	defer query.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	if err = funcs(query); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return nil,nil
}

// this is mysql delete
func (db *MysqlDB) Delete(tablename string,field []string,values ...interface{}) (sql.Result,error) {
	DBSet := db.SqlDB
	strtag,strvalue := parseSliceString(field)
	prepare := fmt.Sprintf("DELETE FROM %v WHERE %v = %v",tablename,strtag,strvalue)
	stmtIns , err := DBSet.Prepare(prepare)
	defer stmtIns.Close()
	if err != nil {
		return nil,fmt.Errorf("db prepare is error : %v",err)
	}
	result, err := stmtIns.Exec(values...)
	if err != nil {
		return nil,fmt.Errorf("db exec is error : %v",err)
	}
	return result,nil
}

// this is mysql update delete
func (db *MysqlDB) UPDATE(tablename string,field []string,screen string,values ...interface{}) (sql.Result,error) {
	DBset := db.SqlDB
	strtag := parseSliceTagAndValue(field)
	prepare := fmt.Sprintf("UPDATE %v SET %v WHERE %v = ?",tablename,strtag,screen)
	stmtIns , err := DBset.Prepare(prepare)
	if err != nil {
		return nil,fmt.Errorf("db prepare is error : %v",err)
	}
	res , err := stmtIns.Exec(values...)
	if err != nil {
		return nil,fmt.Errorf("db exec is error : %v",err)
	}
	return res,nil
}


// this is mysql custom prepare
func (db *MysqlDB) CustomPrepare(sqlflag string,values ...interface{}) (sql.Result,error) {
	DBset := db.SqlDB
	stemIns , err := DBset.Prepare(sqlflag)
	if err != nil {
		return nil,fmt.Errorf("db prepare is error",err)
	}
	res , err := stemIns.Exec(values...)
	if err != nil {
		return nil,fmt.Errorf("db exec is error : %v",err)
	}
	return res,nil

}